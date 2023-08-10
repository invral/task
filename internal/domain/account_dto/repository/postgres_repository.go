package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"task/internal/domain/Errors"
	"task/internal/domain/account_dto/dto"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: pool,
	}
}

func (r *PostgresRepository) UpdateBalance(ctx context.Context, account_id uint64, balance float64) error {
	const op = "PostgresRepository.UpdateBalance"

	_, err := r.CheckExistsAccount(ctx, account_id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query := `
		UPDATE account
			SET balance = @balance,
		WHERE id = @id;
	`

	args := pgx.NamedArgs{
		"id":      account_id,
		"balance": balance,
	}

	if _, err = r.db.Exec(ctx, query, args); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *PostgresRepository) CheckExistsAccount(ctx context.Context, account_id uint64) (*dto.RegistrationCommand, error) {
	const op = "PostgresRepository.CheckExistsAccount"

	query := `
		SELECT id, balance, currency FROM account
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": account_id,
	}

	var accountDto dto.RegistrationCommand

	if err := pgxscan.Get(ctx, r.db, &accountDto, query, args); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, Errors.ErrAccountNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &accountDto, nil
}
