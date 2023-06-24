package stock

type Sku uint32

type Stocks []Stock

type Stock struct {
	WarehouseId WarehouseId
	Count       Count
}

type WarehouseId int64

type Count uint64

type StocksWithSku []StockWithSku

type StockWithSku struct {
	Sku Sku
	Stock
}
