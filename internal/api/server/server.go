package server

import (
	"net/http"
	"task/common"
	"task/internal/api/response"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
	acc "task/internal/domain/account/controller/handler"
	trans "task/internal/domain/transaction/controller/handler"
)

const JSONContentType = "application/json"

type Server struct {
	account     *acc.Handlers
	transaction *trans.Handlers
}

func NewServer(di *common.DependencyContainer) *Server {
	return &Server{
		account:     acc.NewHandlers(di),
		transaction: trans.NewHandlers(di),
	}
}

func (s *Server) GetHTTPHandler(logger *slog.Logger) (http.Handler, error) {

	r := chi.NewRouter()

	r.Use(
		middleware.Recoverer,

		middleware.AllowContentType(JSONContentType),
		render.SetContentType(render.ContentTypeJSON),
		middleware.RequestID,
		common.NewHandler(logger),
		common.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
			requestGroup := slog.Group(
				"request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)

			responseGroup := slog.Group("response",
				slog.Int("status", status),
				slog.Int("bytes", size),
				slog.Duration("duration", duration),
			)

			common.FromRequest(r).InfoCtx(r.Context(), "request processed", requestGroup, responseGroup)
		}),
	)

	r.Group(func(r chi.Router) {
		r.Post("/accounts/register", ErrorHandler(s.account.Register))
		r.Get("/accounts/{account_id}", ErrorHandler(s.account.Get))
		r.Patch("/accounts/{account_id}", ErrorHandler(s.account.Update))
		r.Delete("/accounts/delete", ErrorHandler(s.account.Delete))

		r.Post("/transaction/deposit", ErrorHandler(s.transaction.Deposit))
		r.Post("/transaction/withdraw", ErrorHandler(s.transaction.Withdraw))
		r.Get("/transaction/{transaction_id}", ErrorHandler(s.transaction.GetTransactionByID))
		r.Patch("/transaction/{transaction_id}", ErrorHandler(s.transaction.UpdateTransactionStatus))
		r.Delete("/transaction/{transaction_id}", ErrorHandler(s.transaction.DeleteTransactionByID))
		r.Get("/transaction/frozen/{account_id}", ErrorHandler(s.transaction.GetFrozenBalanceByID))
	})

	return r, nil
}

func ErrorHandler(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			render.JSON(w, r, response.Response{Error: "Internal server error in handlers func", Status: "error"})
		}
	}
}
