package request

import (
	"github.com/go-chi/render"
	"net/http"
	"task/internal/api/response"
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
	ID      uint64  `json:"id"`
	Balance float64 `json:"balance"`
}

func ResponseFrozenBalanceOK(w http.ResponseWriter, r *http.Request, id uint64, balance float64) {
	render.JSON(w, r, ResponseFrozenBalance{
		Response: response.Response{
			Status: response.StatusSuccess,
		},
		ID:      id,
		Balance: balance,
	})
}
