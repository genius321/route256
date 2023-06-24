package schema

type Stock struct {
	WarehouseID int64 `db:"warehouse_id"`
	Amount      int64 `db:"amount"`
}

type StockWithSku struct {
	Sku int64 `db:"sku"`
	Stock
}
