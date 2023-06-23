package postgres

import (
	"context"
	"fmt"
	"log"
	orderModels "route256/loms/internal/models/order"
	"route256/loms/internal/pkg/loms"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *Repository) CreateOrder(
	ctx context.Context,
	user orderModels.User,
	items orderModels.Items,
) (orderModels.OrderId, error) {
	db := r.provider.GetDB(ctx)

	query := psql.Insert(tableNameOrders).Columns("user_id").
		Values(user).
		Suffix("RETURNING order_id")

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query for create order in orders: %s", err)
	}

	var orderId orderModels.OrderId
	err = db.QueryRow(ctx, rawSQL, args...).Scan(&orderId)
	log.Println("NEW ORDER: ", orderId)
	if err != nil {
		return 0, fmt.Errorf("exec insert order in orders: %w", err)
	}

	for _, v := range items {
		query = psql.Insert(tableNameOrderItems).Columns("order_id", "sku", "amount").
			Values(orderId, v.Sku, v.Count)
		rawSQL, args, err = query.ToSql()
		if err != nil {
			return 0, fmt.Errorf("build query for create item in order_items: %s", err)
		}
		_, err = db.Exec(ctx, rawSQL, args...)
		if err != nil {
			return 0, fmt.Errorf("exec insert item in order_items: %w", err)
		}
	}

	return orderId, nil
}

func (r *Repository) ListOrder(ctx context.Context, req *loms.ListOrderRequest) (*loms.ListOrderResponse, error) {
	db := r.provider.GetDB(ctx)
	query := psql.Select("status_name", "user_id").
		From(tableNameOrders).
		Where(sq.Eq{"order_id": req.OrderId})
	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query ListOrder: %s", err)
	}
	var res loms.ListOrderResponse
	err = db.QueryRow(ctx, rawSQL, args...).Scan(&res.Status, &res.User)
	if err != nil {
		return nil, fmt.Errorf("exec query ListOrder: %w", err)
	}

	query = psql.Select("sku", "amount").
		From(tableNameOrderItems).
		Where(sq.Eq{"order_id": req.OrderId})
	rawSQL, args, err = query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query ListOrder: %s", err)
	}
	var resultSQL []struct {
		Sku   int64 `db:"sku"`
		Count int64 `db:"amount"`
	}
	err = pgxscan.Select(ctx, db, &resultSQL, rawSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query ListOrder: %w", err)
	}

	res.Items = make([]*loms.Item, 0, len(resultSQL))
	for _, v := range resultSQL {
		res.Items = append(res.Items, &loms.Item{Sku: uint32(v.Sku), Count: uint32(v.Count)})
	}
	return &res, nil
}

func (r *Repository) OrderPayed(ctx context.Context, req *loms.OrderPayedRequest) (*emptypb.Empty, error) {
	db := r.provider.GetDB(ctx)
	query := psql.Update(tableNameOrders).
		Set("status_name", "payed").
		Where(sq.Eq{"order_id": req.OrderId})
	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build update status_name OrderPayed: %s", err)
	}
	_, err = db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("exec update status_name OrderPayed: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (r *Repository) CancelOrder(ctx context.Context, req *loms.CancelOrderRequest) (*emptypb.Empty, error) {
	db := r.provider.GetDB(ctx)
	query := psql.Update(tableNameOrders).
		Set("status_name", "cancelled").
		Where(sq.Eq{"order_id": req.OrderId})
	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build update status_name CancelOrder: %s", err)
	}
	_, err = db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("exec update status_name CancelOrder: %w", err)
	}
	return &emptypb.Empty{}, nil
}
