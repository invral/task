package dto

type AccountDTO struct {
	ID       uint64  `json:"id"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
	Email    string  `json:"email"`
}

type RegistrationCommand struct {
	ID       uint64  `json:"id"`
	Currency string  `json:"firstName"`
	Balance  float64 `json:"lastName"`
}
