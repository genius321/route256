//go:generate minimock -i Repository -o ./mocks/ -s "_minimock.go" -g
//go:generate minimock -i TransactionManager -o ./mocks/ -s "_minimock.go" -g
//go:generate minimock -i LomsClient -o ./mocks/ -s "_minimock.go" -g
package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"route256/checkout/internal/config"
	"route256/checkout/internal/pkg/checkout"
	"route256/checkout/internal/pkg/loms"
	"route256/checkout/internal/pkg/product-service"
	"route256/checkout/internal/pkg/ratelimit"
	"route256/checkout/internal/pkg/workerpool"
	"route256/checkout/internal/repository/schema"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
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
	ListCart(ctx context.Context, req *checkout.ListCartRequest) ([]*schema.Item, error)
	TakeCountSkuUserFromCart(ctx context.Context, userId int64, sku int64) (int64, error)
	SubFromCart(ctx context.Context, userId int64, sku int64, count int64) error
}

type LomsClient interface {
	CreateOrder(ctx context.Context, in *loms.CreateOrderRequest, opts ...grpc.CallOption) (*loms.CreateOrderResponse, error)
	ListOrder(ctx context.Context, in *loms.ListOrderRequest, opts ...grpc.CallOption) (*loms.ListOrderResponse, error)
	OrderPayed(ctx context.Context, in *loms.OrderPayedRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	CancelOrder(ctx context.Context, in *loms.CancelOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Stocks(ctx context.Context, in *loms.StocksRequest, opts ...grpc.CallOption) (*loms.StocksResponse, error)
}

type Service struct {
	checkout.UnimplementedCheckoutServer
	LomsClient
	productClient product.ProductServiceClient
	TransactionManager
	Repository
	ratelimit.Ratelimit
}

func NewCheckoutServer(lomsClient loms.LomsClient,
	productClient product.ProductServiceClient,
	repo Repository,
	provider TransactionManager,
	ratelimiter *ratelimit.Ratelimit,
) *Service {
	return &Service{
		LomsClient:         lomsClient,
		productClient:      productClient,
		Repository:         repo,
		TransactionManager: provider,
		Ratelimit:          *ratelimiter,
	}
}

const (
	worker = 5
	limit  = 10
)

func (s *Service) AddToCart(ctx context.Context, req *checkout.AddToCartRequest) (*emptypb.Empty, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	stocksResp, _ := s.LomsClient.Stocks(ctx, &loms.StocksRequest{Sku: req.Sku})
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

func (s *Service) DeleteFromCart(ctx context.Context, req *checkout.DeleteFromCartRequest) (*emptypb.Empty, error) {
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

func (s *Service) ListCart(ctx context.Context, req *checkout.ListCartRequest) (*checkout.ListCartResponse, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	items, err := s.Repository.ListCart(ctx, req)
	if err != nil {
		return nil, err
	}
	respItems := make([]*checkout.Item, len(items))

	var totalPrice atomic.Int64

	wp := workerpool.New[product.GetProductRequest, product.GetProductResponse](worker)

	timeStart := time.Now()
	// для поимки ошибки среди работы множества горутин
	g, gCtx := errgroup.WithContext(ctx)
	for i, v := range items {
		// Проверка действителен ли контекст
		if err := gCtx.Err(); err != nil {
			log.Println(time.Since(timeStart))
			return nil, err
		}
		eitherCh := wp.Exec(gCtx,
			&product.GetProductRequest{
				Token: config.Token,
				Sku:   uint32(v.Sku),
			},
			s.productClient.GetProduct,
			&s.Ratelimit,
		)
		i, v := i, v
		g.Go(func() error {
			either := <-eitherCh
			log.Println(either.Value, either.Err)
			if either.Err != nil {
				return either.Err
			}
			// у каждого respItems своя ячейка, конкурентности не будет
			respItems[i] = &checkout.Item{
				Sku:   uint32(v.Sku),
				Count: uint32(v.Amount),
				Name:  either.Value.Name,
				Price: either.Value.Price,
			}
			totalPrice.Add(int64(either.Value.Price * uint32(v.Amount)))
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		log.Println(time.Since(timeStart))
		return nil, err
	}
	log.Println(time.Since(timeStart))
	return &checkout.ListCartResponse{Items: respItems, TotalPrice: uint32(totalPrice.Load())}, nil
}

func (s *Service) Purchase(ctx context.Context, req *checkout.PurchaseRequest) (*checkout.PurchaseResponse, error) {
	log.Printf("%+v", req)
	err := req.ValidateAll()
	if err != nil {
		return nil, err
	}
	items, err := s.Repository.ListCart(ctx, &checkout.ListCartRequest{User: req.User})
	if err != nil {
		return nil, err
	}
	respItems := make([]*loms.Item, 0, len(items))
	for _, v := range items {
		respItems = append(respItems, &loms.Item{
			Sku:   uint32(v.Sku),
			Count: uint32(v.Amount),
		})
	}
	var createOrderResp *loms.CreateOrderResponse
	createOrderResp, err = s.LomsClient.CreateOrder(ctx, &loms.CreateOrderRequest{
		User:  req.User,
		Items: respItems,
	})
	if err != nil {
		return nil, fmt.Errorf("purchase: %w", err)
	}
	_, err = s.Repository.DeleteAllFromCart(ctx, &checkout.DeleteFromCartRequest{User: req.User})

	if err != nil {
		return nil, fmt.Errorf("purchase: %w", err)
	}
	return &checkout.PurchaseResponse{OrderId: createOrderResp.OrderId}, nil
}
