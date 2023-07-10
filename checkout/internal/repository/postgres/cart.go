package postgres

import (
	"context"
	"fmt"
	"route256/checkout/internal/pkg/checkout"
	"route256/checkout/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (r *Repository) AddToCart(ctx context.Context, req *checkout.AddToCartRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "internal/repository/postgres/AddToCart")
	defer span.Finish()
	db := r.provider.GetDB(ctx)
	query := `
INSERT INTO carts("user_id", "sku", "amount") VALUES 
    ($1, $2, $3)
ON CONFLICT ("user_id", "sku") DO UPDATE 
	SET amount=carts.amount+$3
`
	_, err := db.Exec(ctx, query, req.User, req.Sku, req.Count)
	if err != nil {
		return nil, fmt.Errorf("exec insert cart: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (r *Repository) TakeCountSkuUserFromCart(ctx context.Context, userId int64, sku int64) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "internal/repository/postgres/TakeCountSkuUserFromCart")
	defer span.Finish()
	db := r.provider.GetDB(ctx)
	query := psql.Select("amount").
		From("carts").
		Where(sq.Eq{"user_id": userId, "sku": sku})
	rawSQL, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build select TakeCountSkuUserFromCart %s", err)
	}
	var cnt int64
	err = db.QueryRow(ctx, rawSQL, args...).Scan(&cnt)
	if err != nil {
		return 0, fmt.Errorf("exec select TakeCountSkuUserFromCart: %w", err)
	}
	return cnt, nil
}

func (r *Repository) SubFromCart(ctx context.Context, userId int64, sku int64, count int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "internal/repository/postgres/SubFromCart")
	defer span.Finish()
	db := r.provider.GetDB(ctx)
	query := `
UPDATE carts
SET amount=carts.amount-$3
WHERE user_id = $1 and sku = $2;
`
	_, err := db.Exec(ctx, query, userId, sku, count)
	if err != nil {
		return fmt.Errorf("exec update SubFromCart: %w", err)
	}
	return nil
}

func (r *Repository) DeleteFromCart(ctx context.Context, req *checkout.DeleteFromCartRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "internal/repository/postgres/DeleteFromCart")
	defer span.Finish()
	db := r.provider.GetDB(ctx)
	query := `
DELETE FROM carts
WHERE user_id = $1 and sku = $2;
`
	_, err := db.Exec(ctx, query, req.User, req.Sku)
	if err != nil {
		return nil, fmt.Errorf("exec delete DeleteFromCart: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (r *Repository) DeleteAllFromCart(ctx context.Context, req *checkout.DeleteFromCartRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "internal/repository/postgres/DeleteAllFromCart")
	defer span.Finish()
	db := r.provider.GetDB(ctx)
	query := `
DELETE FROM carts
WHERE user_id = $1;
`
	_, err := db.Exec(ctx, query, req.User)
	if err != nil {
		return nil, fmt.Errorf("exec delete DeleteAllFromCart: %w", err)
	}
	return &emptypb.Empty{}, nil
}

func (r *Repository) ListCart(ctx context.Context, req *checkout.ListCartRequest) ([]*schema.Item, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "internal/repository/postgres/ListCart")
	defer span.Finish()
	db := r.provider.GetDB(ctx)
	query := psql.Select("sku", "amount").
		From("carts").
		Where(sq.Eq{"user_id": req.User})
	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select ListCart: %s", err)
	}
	var resultSQL []*schema.Item
	err = pgxscan.Select(ctx, db, &resultSQL, rawSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("exec select ListCart: %w", err)
	}
	return resultSQL, nil
}
