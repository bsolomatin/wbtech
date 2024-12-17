package models

type Payment struct {
	Transaction     string  `json:"transaction"`
	RequestId       string  `json:"request_id"`
	Currency        string  `json:"currency"`
	Provider        string  `json:"provider"`
	Amount          float64 `json:"amount"`
	PaymentDateTime int64   `json:"payment_dt"` // unix time
	Bank            string  `json:"bank"`
	DeliveryCost    float64 `json:"delivery_cost"`
	GoodsTotal      float64 `json:"goods_total"`
	CustomFee       int     `json:"custom_fee"`
}
