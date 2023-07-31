package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"task/internal/domain/Errors"
	"task/internal/domain/account/entity"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: pool,
	}
}

func (r *PostgresRepository) Save(ctx context.Context, entity *entity.Account) error {

	const op = "domain/account.PostgresRepository.Save"
	query := `
		INSERT INTO account (id, currency, balance, password, email)
		VALUES (@id, @currency, @balance, @password, @email)
	`

	args := pgx.NamedArgs{
		"id":       entity.ID,
		"currency": entity.Currency,
		"balance":  entity.Balance,
		"password": entity.Password,
		"email":    entity.Email,
	}

	if _, err := r.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *PostgresRepository) Get(ctx context.Context, id uint64) (*entity.Account, error) {
	const op = "domain/account.PostgresRepository.Get"
	query := `
		SELECT id, currency, balance, password, email FROM account
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	var account entity.Account
	var err error
	if err = pgxscan.Get(ctx, r.db, &account, query, args); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, Errors.ErrAccountNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &account, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uint64) error {
	const op = "domain/account.PostgresRepository.Delete"
	query := `
		DELETE FROM account
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	if _, err := r.Get(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, Errors.ErrAccountNotFound)
	}

	if _, err := r.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *PostgresRepository) Update(ctx context.Context, id uint64, balance float64, currency string) error {
	const op = "domain/account.PostgresRepository.Update"
	query := `
		UPDATE account
			SET balance = @balance,
				currency = @currency
		WHERE id = @id;
	`

	args := pgx.NamedArgs{
		"id":       id,
		"balance":  balance,
		"currency": currency,
	}

	if _, err := r.Get(ctx, id); err != nil {
		return fmt.Errorf("%s: %w", op, Errors.ErrAccountNotFound)
	}

	if _, err := r.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
