package models

type Payload struct {
	Order Order
	Delivery Delivery
	Payment Payment
	Items []Item
}