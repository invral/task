package entity

type Account struct {
	ID       uint64  `json:"id"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
	Password string  `json:"password"`
	Email    string  `json:"email"`
}

func NewAccount(id uint64, currency string, balance float64, password string, email string) *Account {
	return &Account{
		ID:       id,
		Currency: currency,
		Balance:  balance,
		Password: password,
		Email:    email,
	}
}
