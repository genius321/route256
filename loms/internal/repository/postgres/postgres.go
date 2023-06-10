package postgres

import (
	"route256/loms/internal/repository/postgres/tx"

	sq "github.com/Masterminds/squirrel"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repository struct {
	provider tx.DBProvider
}

func New(provider tx.DBProvider) *Repository {
	return &Repository{provider: provider}
}

const (
	tableNameOrders               = "orders"
	tableNameOrderItems           = "order_items"
	tableNameSkuStocks            = "sku_stocks"
	tableNameSkuStocksReservation = "sku_stocks_reservation"
)
