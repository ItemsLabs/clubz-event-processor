package types

type PurchaseDetails struct {
	ItemId        string  `json:"item_id"`
	Quantity      int64   `json:"quantity"`
	Amount        float64 `json:"amount"`
	PaymentMethod int64   `json:"payment_method"` // 1 - stripe or 2 - internal
	CancelUrl     string  `json:"cancel_url"`
	SuccessUrl    string  `json:"success_url"`
}
