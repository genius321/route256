package postgres

import (
	"context"
	"fmt"
	stockModels "route256/loms/internal/models/stock"
	"route256/loms/internal/repository/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

func (r *Repository) Stocks(ctx context.Context, sku stockModels.Sku) (stockModels.Stocks, error) {
	db := r.provider.GetDB(ctx)

	query := psql.Select("warehouse_id", "amount").
		From(tableNameSkuStocks).
		Where(sq.Eq{"sku": sku})

	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query for get stocks: %s", err)
	}

	// думал вынести в модель на этом слое, но не придумал название
	var resultSQL []struct {
		WarehouseID int64 `db:"warehouse_id"`
		Count       int64 `db:"amount"`
	}

	err = pgxscan.Select(ctx, db, &resultSQL, rawSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query get stocks: %s", err)
	}

	stocks := make(stockModels.Stocks, 0, len(resultSQL))
	for _, v := range resultSQL {
		stocks = append(stocks, stockModels.Stock{
			WarehouseId: stockModels.WarehouseId(v.WarehouseID),
			Count:       stockModels.Count(v.Count),
		})
	}

	return stocks, nil
}

func (r *Repository) AddSkuStockReserve(ctx context.Context, sku int64, amount int64, warehouseID int64, orderId int64) error {
	db := r.provider.GetDB(ctx)
	query := `
INSERT INTO sku_stocks_reservation("sku", "warehouse_id", "order_id", "amount") VALUES 
    ($1, $2, $3, $4)
`
	_, err := db.Exec(ctx, query, sku, warehouseID, orderId, amount)
	if err != nil {
		return fmt.Errorf("exec insert reservation: %w", err)
	}

	return nil
}

func (r *Repository) DeleteStocksReserveByOrderId(ctx context.Context, orderId int64) error {
	db := r.provider.GetDB(ctx)
	query := psql.Delete(tableNameSkuStocksReservation).
		Where(sq.Eq{"order_id": orderId})
	rawSQL, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("build query for delete stocks reserve by orderId: %s", err)
	}
	_, err = db.Exec(ctx, rawSQL, args...)
	if err != nil {
		return fmt.Errorf("exec query for delete stocks reserve by orderId: %w", err)
	}
	return nil
}

func (r *Repository) TakeStocksReserveByOrderId(ctx context.Context, orderId int64) ([]schema.Stocks, error) {
	db := r.provider.GetDB(ctx)
	query := psql.Select("sku", "warehouse_id", "amount").
		From(tableNameSkuStocksReservation).
		Where(sq.Eq{"order_id": orderId})
	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query for select stocks reserve by orderId: %s", err)
	}
	var resultSQL []schema.Stocks

	err = pgxscan.Select(ctx, db, &resultSQL, rawSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query for select stocks reserve by orderId: %s", err)
	}
	return resultSQL, nil
}

func (r *Repository) AddSkuStock(ctx context.Context, sku int64, amount int64, warehouseID int64) error {
	db := r.provider.GetDB(ctx)
	query := `
INSERT INTO sku_stocks("sku", "warehouse_id", "amount") VALUES 
    ($1, $2, $3)
ON CONFLICT ("sku", "warehouse_id") DO UPDATE 
	SET amount=sku_stocks.amount+$3
`

	_, err := db.Exec(ctx, query, sku, warehouseID, amount)
	if err != nil {
		return fmt.Errorf("exec insert stocks: %w", err)
	}

	return nil
}

func (r *Repository) TakeSkuStock(ctx context.Context, sku int64, amount int64, warehouseID int64) (int64, error) {
	db := r.provider.GetDB(ctx)
	query := `
INSERT INTO sku_stocks("sku", "warehouse_id", "amount") VALUES 
    ($1, $2, $3)
ON CONFLICT ("sku", "warehouse_id") DO UPDATE 
	SET amount=sku_stocks.amount-$3
RETURNING amount;
`

	var cnt int64
	err := db.QueryRow(ctx, query, sku, warehouseID, amount).Scan(&cnt)
	if err != nil {
		return 0, fmt.Errorf("exec insert stocks: %w", err)
	}

	return cnt, nil
}
