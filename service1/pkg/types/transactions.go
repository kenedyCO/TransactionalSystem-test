package types

type Transactions struct {
	ID           string `json:"id"`
	IDSender     string `json:"id_sender"`
	Wallet       string `json:"wallet"`
	CurrencyCode string `json:"currencyCode"`
	Amount       int    `json:"amount"`
}
