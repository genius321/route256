package schema

type Item struct {
	Sku    int64 `db:"sku"`
	Amount int64 `db:"amount"`
}
