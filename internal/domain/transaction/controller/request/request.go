package request

import (
	"github.com/go-chi/render"
	"net/http"
	"task/internal/api/response"
	"task/internal/domain/account_dto/dto"
	"task/internal/domain/transaction/entity"
)

type ResponseTransaction struct {
	response.Response
	entity.Transaction
}

func ResponseTransactionOK(w http.ResponseWriter, r *http.Request, transaction entity.Transaction) {
	render.JSON(w, r, ResponseTransaction{
		Response: response.Response{
			Status: response.StatusSuccess,
		},
		Transaction: transaction,
	})
}

func ResponseOK(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, response.Response{
		Status: response.StatusSuccess,
	})
}

type ResponseFrozenBalance struct {
	response.Response
	dto.RegistrationCommand
}

func ResponseFrozenBalanceOK(w http.ResponseWriter, r *http.Request, dto *dto.RegistrationCommand) {
	render.JSON(w, r, ResponseFrozenBalance{
		Response: response.Response{
			Status: response.StatusSuccess,
		},
		RegistrationCommand: *dto,
	})
}
