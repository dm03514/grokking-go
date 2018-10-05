package deposits

type Deposit struct {
	ID            int `json:"id,omitempty"`
	TransactionID int `json:"transaction_id"`
	AccountNumber int `json:"account_number"`
	AmountCents   int `json:"amount_cents"`
}
