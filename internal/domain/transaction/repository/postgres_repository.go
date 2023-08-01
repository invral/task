package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"task/internal/domain/transaction/entity"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
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
			currency
		) VALUES (
			@id,
			@status,
			@account_id,
			@amount,
			@currency
		)`

	args := pgx.NamedArgs{
		"id":         transaction.ID,
		"status":     transaction.Status,
		"account_id": transaction.AccountID,
		"amount":     transaction.Amount,
		"currency":   transaction.Currency,
	}

	if _, err := r.pool.Exec(ctx, query, args); err != nil {
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

	if _, err := r.pool.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
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

	if _, err := r.pool.Exec(ctx, query, args); err != nil {
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

	if _, err := r.pool.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uint64) error {
	const op = "PostgresRepository.Delete"

	query := `
		DELETE FROM transactions
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	if _, err := r.pool.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
