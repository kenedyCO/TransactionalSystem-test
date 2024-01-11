package types

type Transactions struct {
	ID           string  `json:"id"`
	CurrencyCode string  `json:"currencyCode"`
	Amount       float64 `json:"amount"`
	Wallet       string  `json:"wallet"`
}
