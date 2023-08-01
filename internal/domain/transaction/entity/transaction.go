package entity

type Transaction struct {
	ID        uint64  `json:"id,omitempty"`
	Status    string  `json:"status,omitempty"`
	AccountID uint64  `json:"account_id,omitempty"`
	Amount    float64 `json:"amount,omitempty"`
	Currency  string  `json:"currency,omitempty"`
	ToAccount uint64  `json:"to_account,omitempty"`
}
