package handler

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"io"
	"net/http"
	"strconv"
	"task/common"
	"task/internal/api/response"
	"task/internal/domain/transaction/controller/request"

	//"task/internal/domain/account/controller/handler/request"
	"task/internal/domain/transaction/entity"
	"task/internal/domain/transaction/service"
)

type Handlers struct {
	service *service.Service
}

func NewHandlers(di *common.DependencyContainer) *Handlers {
	return &Handlers{
		service: service.NewService(di),
	}
}

func (h *Handlers) Deposit(w http.ResponseWriter, r *http.Request) error {
	const op = "transaction.Handlers.Deposit"
	ctx := r.Context()

	var transaction entity.Transaction

	err := render.DecodeJSON(r.Body, &transaction)

	if errors.Is(err, io.EOF) {
		render.JSON(w, r, response.Response{Error: "empty request", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	if err != nil {
		render.JSON(w, r, response.Response{Error: "failed to decode request", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := h.service.CreateDepositTransaction(ctx, &transaction); err != nil {
		render.JSON(w, r, response.Response{Error: "failed to save transaction", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	request.ResponseTransactionOK(w, r, transaction)

	return nil
}

func (h *Handlers) Withdraw(w http.ResponseWriter, r *http.Request) error {
	const op = "transaction.Handlers.Withdraw"
	ctx := r.Context()

	var transaction entity.Transaction

	err := render.DecodeJSON(r.Body, &transaction)

	if errors.Is(err, io.EOF) {
		render.JSON(w, r, response.Response{Error: "empty request", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	if err != nil {
		render.JSON(w, r, response.Response{Error: "failed to decode request", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := h.service.CreateWithdrawTransaction(ctx, &transaction); err != nil {
		render.JSON(w, r, response.Response{Error: "failed to save transaction", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	request.ResponseTransactionOK(w, r, transaction)

	return nil
}

func (h *Handlers) GetTransactionByID(w http.ResponseWriter, r *http.Request) error {
	const op = "transaction.Handlers.GetTransactionByID"
	ctx := r.Context()

	id, err := GetIDFromRequest(r, "transaction_id")
	if err != nil {
		render.JSON(w, r, response.Response{Error: "failed to decode request", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	transaction, err := h.service.GetTransactionByID(ctx, id)
	if err != nil {
		render.JSON(w, r, response.Response{Error: "failed to get account", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	request.ResponseTransactionOK(w, r, *transaction)

	return nil
}

func (h *Handlers) UpdateTransactionStatus(w http.ResponseWriter, r *http.Request) error {
	const op = "transaction.Handlers.UpdateTransactionStatus"
	ctx := r.Context()

	var id uint64

	err := render.DecodeJSON(r.Body, &id)

	if errors.Is(err, io.EOF) {
		render.JSON(w, r, response.Response{Error: "empty request", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	if err != nil {
		render.JSON(w, r, response.Response{Error: "failed to decode request", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	err = h.service.UpdateTransactionStatus(ctx, id)
	if err != nil {
		render.JSON(w, r, response.Response{Error: "failed to update transaction status", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	request.ResponseOK(w, r)

	return nil
}

func (h *Handlers) DeleteTransactionByID(w http.ResponseWriter, r *http.Request) error {
	const op = "transaction.Handlers.DeleteTransactionByID"
	ctx := r.Context()

	id, err := GetIDFromRequest(r, "transaction_id")
	if err != nil {
		render.JSON(w, r, response.Response{Error: "failed to decode request", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	err = h.service.DeleteTransactionByID(ctx, id)
	if err != nil {
		render.JSON(w, r, response.Response{Error: "failed to delete transaction", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	request.ResponseOK(w, r)

	return nil
}

func (h *Handlers) GetFrozenBalanceByID(w http.ResponseWriter, r *http.Request) error {
	const op = "transaction.Handlers.GetTransactions"
	ctx := r.Context()

	var id uint64

	err := render.DecodeJSON(r.Body, &id)

	if errors.Is(err, io.EOF) {
		render.JSON(w, r, response.Response{Error: "empty request", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	if err != nil {
		render.JSON(w, r, response.Response{Error: "failed to decode request", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	balance, err := h.service.GetFrozenBalanceByAccountID(ctx, id)
	if err != nil {
		render.JSON(w, r, response.Response{Error: "failed to get frozen balance", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	request.ResponseFrozenBalanceOK(w, r, id, balance)

	return nil
}

func GetIDFromRequest(r *http.Request, key string) (uint64, error) {
	param := chi.URLParam(r, key)
	if param == "" {
		return 0, fmt.Errorf("empty parameter %s", key)
	}

	value, err := strconv.ParseUint(param, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid parameter %s", key)
	}

	return value, nil
}
