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
	"task/internal/domain/account/controller/handler/request"
	"task/internal/domain/account/entity"
	"task/internal/domain/account/service"
)

const (
	USD = "USD"
	EUR = "EUR"
	RUB = "RUB"
)

type Handlers struct {
	service *service.Service
}

func NewHandlers(di *common.DependencyContainer) *Handlers {
	return &Handlers{
		service: service.NewService(di),
	}
}

func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) error {
	const op = "account.Handlers.Register"
	ctx := r.Context()

	var account entity.Account

	err := render.DecodeJSON(r.Body, &account)

	if errors.Is(err, io.EOF) {
		render.JSON(w, r, response.Response{Error: "empty request", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	if err != nil {
		render.JSON(w, r, response.Response{Error: "failed to decode request", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	if account.Currency != USD && account.Currency != EUR && account.Currency != RUB {
		render.JSON(w, r, response.Response{Error: "invalid currency", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err := h.service.SaveAccount(ctx, &account); err != nil {
		render.JSON(w, r, response.Response{Error: "failed to save account", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	request.ResponseRegisterOK(w, r, &account)

	return nil

}

func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) error {
	const op = "account.Handlers.Get"
	ctx := r.Context()

	id, err := GetIDFromRequest(r, "account_id")
	if err != nil {
		render.JSON(w, r, response.Response{Error: "failed to decode request", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	account, err := h.service.GetAccount(ctx, id)
	if err != nil {
		render.JSON(w, r, response.Response{Error: "failed to get account", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	request.ResponseGetOK(w, r, account)

	return nil

}

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) error {
	const op = "account.Handlers.Delete"
	ctx := r.Context()

	id, err := GetIDFromRequest(r, "account_id")
	if err != nil {
		render.JSON(w, r, response.Response{Error: "failed to get id", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	err = h.service.DeleteAccount(ctx, id)
	if err != nil {
		render.JSON(w, r, response.Response{Error: "failed to delete account", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	render.JSON(w, r, response.Response{
		Status: "ok",
	})

	return nil
}

func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) error {
	const op = "account.Handlers.Update"
	ctx := r.Context()

	var req request.Request

	err := render.DecodeJSON(r.Body, &req)

	if errors.Is(err, io.EOF) {
		render.JSON(w, r, response.Response{Error: "empty request", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	if err != nil {
		render.JSON(w, r, response.Response{Error: "failed to decode", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	if req.Currency != USD && req.Currency != EUR && req.Currency != RUB {
		render.JSON(w, r, response.Response{Error: "invalid currency", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	err = h.service.UpdateBalance(ctx, req.ID, req.Balance, req.Currency)
	if err != nil {
		render.JSON(w, r, response.Response{Error: "failed to update balance", Status: "error"})
		return fmt.Errorf("%s: %w", op, err)
	}

	request.ResponseUpdateOK(w, r, req.ID, req.Balance, req.Currency)

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
