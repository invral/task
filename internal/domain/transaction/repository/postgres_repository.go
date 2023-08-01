package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"task/internal/domain/Errors"
	"task/internal/domain/transaction/entity"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: pool,
	}
}

func (r *PostgresRepository) CreateDepositTransaction(ctx context.Context, transaction *entity.Transaction) error {
	const op = "PostgresRepository.CreateDepositTransaction"

	query := `
		INSERT INTO transactions (
			id,
			status,
			account_id,
			amount,	
			currency,
			to_account
		) VALUES (
			@id,
			@status,
			@account_id,
			@amount,
			@currency,
			NULL
		)`

	args := pgx.NamedArgs{
		"id":         transaction.ID,
		"status":     transaction.Status,
		"account_id": transaction.AccountID,
		"amount":     transaction.Amount,
		"currency":   transaction.Currency,
	}

	if _, err := r.GetTransactionByID(ctx, transaction.ID); err == nil {
		return fmt.Errorf("%s: %w", op, Errors.ErrTransactionExists)
	}

	if _, err := r.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *PostgresRepository) CreateWithdrawTransaction(ctx context.Context, transaction *entity.Transaction) error {
	const op = "PostgresRepository.CreateWithdrawTransaction"

	query := `
		INSERT INTO transactions (
			id,
			status,
			account_id,
			amount,	
			currency,
			to_account

		) VALUES (
			@id,
			@status,
			@account_id,
			@amount,
			@currency
			@to_account
		)`

	args := pgx.NamedArgs{
		"id":         transaction.ID,
		"status":     transaction.Status,
		"account_id": transaction.AccountID,
		"amount":     transaction.Amount,
		"currency":   transaction.Currency,
		"to_account": transaction.ToAccount,
	}

	if _, err := r.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *PostgresRepository) GetTransactionByID(ctx context.Context, id uint64) (*entity.Transaction, error) {
	const op = "transaction.PostgresRepository.GetTransaction"

	query := `
	SELECT id, status, account_id, amount, currency, to_account FROM transactions
	WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	var transaction entity.Transaction

	if err := pgxscan.Get(ctx, r.db, &transaction, query, args); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, Errors.ErrTransactionNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &transaction, nil

}

func (r *PostgresRepository) UpdateWithSuccess(ctx context.Context, id uint64) error {
	const op = "PostgresRepository.UpdateWithSuccess"

	query := `
		UPDATE transactions
		SET status = 'success'
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	if _, err := r.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *PostgresRepository) UpdateWithError(ctx context.Context, id uint64) error {
	const op = "PostgresRepository.UpdateWithError"

	query := `
		UPDATE transactions
		SET status = 'error'
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	if _, err := r.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *PostgresRepository) DeleteTransactionByID(ctx context.Context, id uint64) error {
	const op = "PostgresRepository.Delete"

	query := `
		DELETE FROM transactions
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	if _, err := r.GetTransactionByID(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, Errors.ErrTransactionNotFound)
	}

	if _, err := r.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
