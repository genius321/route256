package stock

type Sku uint32

type Stocks []Stock

type Stock struct {
	WarehouseId WarehouseId
	Count       Count
}

type WarehouseId int64

type Count uint64
