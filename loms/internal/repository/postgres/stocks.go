package postgres

import (
	"context"
	"fmt"
	orderModels "route256/loms/internal/models/order"
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

	var resultSQL []schema.Stock

	err = pgxscan.Select(ctx, db, &resultSQL, rawSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query get stocks: %s", err)
	}

	stocks := make(stockModels.Stocks, 0, len(resultSQL))
	for _, v := range resultSQL {
		stocks = append(stocks, stockModels.Stock{
			WarehouseId: stockModels.WarehouseId(v.WarehouseID),
			Count:       stockModels.Count(v.Amount),
		})
	}

	return stocks, nil
}

func (r *Repository) AddSkuStockReserve(
	ctx context.Context, stockWithSku stockModels.StockWithSku, orderId orderModels.OrderId) error {
	db := r.provider.GetDB(ctx)
	query := `
INSERT INTO sku_stocks_reservation("sku", "warehouse_id", "order_id", "amount") VALUES 
    ($1, $2, $3, $4)
`
	_, err := db.Exec(ctx, query, stockWithSku.Sku, stockWithSku.WarehouseId, orderId, stockWithSku.Count)
	if err != nil {
		return fmt.Errorf("exec insert reservation: %w", err)
	}

	return nil
}

func (r *Repository) DeleteStocksReserveByOrderId(ctx context.Context, orderId orderModels.OrderId) error {
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

func (r *Repository) TakeStocksReserveByOrderId(
	ctx context.Context,
	orderId orderModels.OrderId,
) (stockModels.StocksWithSku, error) {
	db := r.provider.GetDB(ctx)
	query := psql.Select("sku", "warehouse_id", "amount").
		From(tableNameSkuStocksReservation).
		Where(sq.Eq{"order_id": orderId})
	rawSQL, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query for select stocks reserve by orderId: %s", err)
	}

	var resultSQL []schema.StockWithSku

	err = pgxscan.Select(ctx, db, &resultSQL, rawSQL, args...)
	if err != nil {
		return nil, fmt.Errorf("exec query for select stocks reserve by orderId: %s", err)
	}

	res := make(stockModels.StocksWithSku, len(resultSQL))
	for i, v := range resultSQL {
		res[i].Sku = stockModels.Sku(v.Sku)
		res[i].WarehouseId = stockModels.WarehouseId(v.WarehouseID)
		res[i].Count = stockModels.Count(v.Amount)
	}
	return res, nil
}

func (r *Repository) AddSkuStock(ctx context.Context, stockWithSku stockModels.StockWithSku) error {
	db := r.provider.GetDB(ctx)
	query := `
INSERT INTO sku_stocks("sku", "warehouse_id", "amount") VALUES 
    ($1, $2, $3)
ON CONFLICT ("sku", "warehouse_id") DO UPDATE 
	SET amount=sku_stocks.amount+$3
`

	_, err := db.Exec(ctx, query, stockWithSku.Sku, stockWithSku.WarehouseId, stockWithSku.Count)
	if err != nil {
		return fmt.Errorf("exec insert stocks: %w", err)
	}

	return nil
}

func (r *Repository) TakeSkuStock(
	ctx context.Context, stockWithSku stockModels.StockWithSku) (stockModels.Count, error) {
	db := r.provider.GetDB(ctx)
	query := `
INSERT INTO sku_stocks("sku", "warehouse_id", "amount") VALUES 
    ($1, $2, $3)
ON CONFLICT ("sku", "warehouse_id") DO UPDATE 
	SET amount=sku_stocks.amount-$3
RETURNING amount;
`

	var cnt stockModels.Count
	err := db.QueryRow(ctx, query, stockWithSku.Sku, stockWithSku.WarehouseId, stockWithSku.Count).Scan(&cnt)
	if err != nil {
		return 0, fmt.Errorf("exec insert stocks: %w", err)
	}

	return cnt, nil
}
