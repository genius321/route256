package domain

type Item struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

type Stock struct {
	WarehouseID int64
	Count       uint64
}

type Product struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type ItemCart struct {
	Item
	Product
}
