package request

import (
	"github.com/go-chi/render"
	"net/http"
	"task/internal/api/response"
	"task/internal/domain/account/entity"
	"task/internal/domain/account_dto/dto"
)

type Request struct {
	ID       uint64  `json:"id,omitempty"`
	Balance  float64 `json:"balance,omitempty"`
	Currency string  `json:"currency,omitempty"`
	Password string  `json:"password,omitempty" validate:"required,alphanumeric"`
	Email    string  `json:"email,omitempty" validate:"required,email"`
}

type ResponseSave struct {
	response.Response
	dto.RegistrationCommand
}

type ResponseGet struct {
	response.Response
	dto.AccountDTO
}

type ResponseUpdate struct {
	response.Response
	id       uint64
	balance  float64
	currency string
}

func ResponseRegisterOK(w http.ResponseWriter, r *http.Request, account *entity.Account) {
	render.JSON(w, r, ResponseSave{
		Response: response.Response{
			Status: "ok",
		},
		RegistrationCommand: dto.RegistrationCommand{
			ID:       account.ID,
			Balance:  account.Balance,
			Currency: account.Currency,
		},
	})
}

func ResponseGetOK(w http.ResponseWriter, r *http.Request, account *entity.Account) {
	render.JSON(w, r, ResponseGet{
		Response: response.Response{
			Status: "ok",
		},
		AccountDTO: dto.AccountDTO{
			ID:       account.ID,
			Currency: account.Currency,
			Balance:  account.Balance,
			Email:    account.Email,
		},
	})
}

func ResponseUpdateOK(w http.ResponseWriter, r *http.Request, id uint64, balance float64, currency string) {
	render.JSON(w, r, ResponseUpdate{
		Response: response.Response{
			Status: "ok",
		},
		id:       id,
		balance:  balance,
		currency: currency,
	})
}
