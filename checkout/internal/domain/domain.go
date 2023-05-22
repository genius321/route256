package domain

type StockChecker interface {
	Stocks(sku uint32) ([]Stock, error)
}

type Stock struct {
	WarehouseID int64
	Count       uint64
}

type Model struct {
	stockChecker StockChecker
}

func New(stockChecker StockChecker) *Model {
	return &Model{stockChecker: stockChecker}
}
