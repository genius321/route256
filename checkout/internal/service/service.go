package service

import (
	"context"
	"errors"
	"log"
	"route256/checkout/internal/config"
	"route256/checkout/internal/pkg/checkout"
	"route256/checkout/internal/pkg/loms"
	"route256/checkout/internal/pkg/product-service"
	"route256/checkout/internal/pkg/workerpool"
	"route256/checkout/internal/repository/schema"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aitsvet/debugcharts"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TransactionManager interface {
	RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) error
	Serializable(ctx context.Context, fn func(ctxTx context.Context) error) error
}

type Repository interface {
	AddToCart(ctx context.Context, req *checkout.AddToCartRequest) (*emptypb.Empty, error)
	DeleteFromCart(ctx context.Context, req *checkout.DeleteFromCartRequest) (*emptypb.Empty, error)
	DeleteAllFromCart(ctx context.Context, req *checkout.DeleteFromCartRequest) (*emptypb.Empty, error)
	ListCart(ctx context.Context, req *checkout.ListCartRequest) (*[]schema.Item, error)
	TakeCountSkuUserFromCart(ctx context.Context, userId int64, sku int64) (int64, error)
	SubFromCart(ctx context.Context, userId int64, sku int64, count int64) error
}

type service struct {
	checkout.UnimplementedCheckoutServer
	lomsClient    loms.LomsClient
	productClient product.ProductServiceClient
	TransactionManager
	Repository
}

func NewCheckoutServer(lomsClient loms.LomsClient,
	productClient product.ProductServiceClient,
	repo Repository,
	provider TransactionManager,
) *service {
	return &service{
		lomsClient:         lomsClient,
		productClient:      productClient,
		Repository:         repo,
		TransactionManager: provider,
	}
}

const worker = 5

func (s *service) AddToCart(ctx context.Context, req *checkout.AddToCartRequest) (*emptypb.Empty, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	stocksResp, _ := s.lomsClient.Stocks(ctx, &loms.StocksRequest{Sku: req.Sku})
	stocks := stocksResp.Stocks
	log.Printf("stocks: %v", stocks)
	counter := int64(req.Count)
	for _, stock := range stocks {
		counter -= int64(stock.Count)
		if counter <= 0 {
			return s.Repository.AddToCart(ctx, req)
		}
	}
	return nil, errors.New("stock insufficient")
}

func (s *service) DeleteFromCart(ctx context.Context, req *checkout.DeleteFromCartRequest) (*emptypb.Empty, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	count, err := s.Repository.TakeCountSkuUserFromCart(ctx, req.User, int64(req.Sku))
	if err != nil {
		return nil, err
	}
	if count-int64(req.Count) > 0 {
		return &emptypb.Empty{}, s.Repository.SubFromCart(ctx, req.User, int64(req.Sku), int64(req.Count))
	}
	return s.Repository.DeleteFromCart(ctx, req)
}

func (s *service) ListCart(ctx context.Context, req *checkout.ListCartRequest) (*checkout.ListCartResponse, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	items, err := s.Repository.ListCart(ctx, req)
	if err != nil {
		return nil, err
	}
	respItems := make([]*checkout.Item, len(*items))

	var totalPrice atomic.Int64

	wp := workerpool.New[product.GetProductRequest, product.GetProductResponse](worker)

	wg := sync.WaitGroup{}
	timeStart := time.Now()
	for i, v := range *items {
		wg.Add(1)
		debugcharts.RPS.Add(1)
		eitherCh := wp.Exec(ctx,
			&product.GetProductRequest{
				Token: config.Token,
				Sku:   uint32(v.Sku),
			},
			s.productClient.GetProduct,
		)
		go func(i int, v schema.Item) {
			defer wg.Done()
			either := <-eitherCh
			log.Println(either.Value)
			if either.Err != nil {
				log.Fatalf("get product: %v\n", either.Err)
			}
			respItems[i] = &checkout.Item{
				Sku:   uint32(v.Sku),
				Count: uint32(v.Amount),
				Name:  either.Value.Name,
				Price: either.Value.Price,
			}
			totalPrice.Add(int64(either.Value.Price * uint32(v.Amount)))
		}(i, v)
	}
	wg.Wait()
	log.Println(time.Since(timeStart))
	return &checkout.ListCartResponse{Items: respItems, TotalPrice: uint32(totalPrice.Load())}, nil
}

func (s *service) Purchase(ctx context.Context, req *checkout.PurchaseRequest) (*checkout.PurchaseResponse, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	items, err := s.Repository.ListCart(ctx, &checkout.ListCartRequest{User: req.User})
	if err != nil {
		return nil, err
	}
	respItems := make([]*loms.Item, 0, len(*items))
	for _, v := range *items {
		respItems = append(respItems, &loms.Item{
			Sku:   uint32(v.Sku),
			Count: uint32(v.Amount),
		})
	}
	var createOrderResp *loms.CreateOrderResponse
	err = s.Serializable(ctx, func(ctxTx context.Context) error {
		createOrderResp, err = s.lomsClient.CreateOrder(ctx, &loms.CreateOrderRequest{
			User:  req.User,
			Items: respItems,
		})
		if err != nil {
			return err
		}
		_, err := s.Repository.DeleteAllFromCart(ctx, &checkout.DeleteFromCartRequest{User: req.User})
		return err
	})
	return &checkout.PurchaseResponse{OrderId: createOrderResp.OrderId}, nil
}
