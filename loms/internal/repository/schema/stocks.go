package schema

type Stocks struct {
	Sku         int64 `db:"sku"`
	WarehouseID int64 `db:"warehouse_id"`
	Amount      int64 `db:"amount"`
}
