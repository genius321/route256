package order

type User int64

type Items []Item

type Item struct {
	Sku   Sku
	Count Count
}

type Sku uint32

type Count uint16

type OrderId int64

type Status string
