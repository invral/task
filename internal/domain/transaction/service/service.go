package service

import (
	"context"
	"errors"
	"fmt"
	"task/common"
	"task/internal/api/response"
	"task/internal/domain/Errors"
	"task/internal/domain/account_dto/dto"
	rep "task/internal/domain/account_dto/repository"
	"task/internal/domain/transaction/entity"
	"task/internal/domain/transaction/repository"
)

//go:generate go run github.com/vektra/mockery/v2@v2.32.4 --name=Repository_transaction
type Repository_transaction interface {
	CreateDepositTransaction(ctx context.Context, transaction *entity.Transaction) error
	CreateWithdrawTransaction(ctx context.Context, transaction *entity.Transaction) error
	GetTransactionByID(ctx context.Context, id uint64) (*entity.Transaction, error)
	UpdateTransactionStatus(ctx context.Context, id uint64, status string) error
	DeleteTransactionByID(ctx context.Context, id uint64) error
	GetTransactionsByAccountID(ctx context.Context, accountID uint64) ([]*entity.Transaction, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.32.4 --name=Repository_acc_dto
type Repository_acc_dto interface {
	UpdateBalance(ctx context.Context, account_id uint64, balance float64) error
	CheckExistsAccount(ctx context.Context, account_id uint64) (*dto.RegistrationCommand, error)
}

type Service struct {
	repTransaction Repository_transaction
	repAccDto      Repository_acc_dto
}

func NewService(di *common.DependencyContainer) *Service {
	return &Service{
		repTransaction: repository.NewPostgresRepository(di.Pool),
		repAccDto:      rep.NewPostgresRepository(di.Pool),
	}
}

func (s *Service) CreateDepositTransaction(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error) {
	const op = "domain/transaction.Service.CreateDepositTransaction"

	_, err := s.repTransaction.GetTransactionByID(ctx, transaction.ID)

	switch err {
	case nil:
		return nil, fmt.Errorf("%s: %w", op, Errors.ErrTransactionExists)

	case Errors.ErrTransactionNotFound:
		err = s.repTransaction.CreateDepositTransaction(ctx, transaction)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return transaction, nil
	default:
		return nil, fmt.Errorf("%s: %w", op, err)
	}
}

func (s *Service) CreateWithdrawTransaction(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error) {
	const op = "domain/transaction.Service.CreateWithdrawTransaction"

	_, err := s.repTransaction.GetTransactionByID(ctx, transaction.ID)

	switch err {
	case nil:
		return nil, fmt.Errorf("%s: %w", op, Errors.ErrAccountExists)

	case Errors.ErrTransactionNotFound:
		err = s.repTransaction.CreateWithdrawTransaction(ctx, transaction)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return transaction, nil
	default:
		return nil, fmt.Errorf("%s: %w", op, err)
	}
}

func (s *Service) GetTransactionByID(ctx context.Context, id uint64) (*entity.Transaction, error) {
	const op = "domain/transaction.Service.GetTransactionByID"

	transaction, err := s.repTransaction.GetTransactionByID(ctx, id)
	if errors.Is(err, Errors.ErrTransactionNotFound) {
		return nil, fmt.Errorf("%s: %w", op, Errors.ErrTransactionNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return transaction, nil
}

func (s *Service) GetFrozenBalanceByAccountID(ctx context.Context, accountID uint64) (*dto.RegistrationCommand, error) {
	const op = "domain/transaction.Service.GetTransactionsByAccountID"

	transactions, err := s.repTransaction.GetTransactionsByAccountID(ctx, accountID)
	if errors.Is(err, Errors.ErrTransactionNotFound) {
		return nil, fmt.Errorf("%s: %w", op, Errors.ErrTransactionNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	accountDto, err := s.repAccDto.CheckExistsAccount(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var total float64

	for _, transaction := range transactions {
		if transaction.Status == response.StatusCreated {
			amount, err := common.ValidateCurrency(transaction.Currency, accountDto.Currency, transaction.Amount)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			switch {
			case transaction.ToAccount == 0:
				total += amount

			case transaction.ToAccount > 0:
				total -= amount
			default:
				return nil, fmt.Errorf("%s: %w", op, err)
			}

		}
	}
	accountDto.Balance = total

	return accountDto, nil
}

func (s *Service) UpdateTransactionStatus(ctx context.Context, id uint64) error {
	const op = "domain/transaction.Service.UpdateTransactionStatus"

	transaction, err := s.repTransaction.GetTransactionByID(ctx, id)
	if errors.Is(err, Errors.ErrTransactionNotFound) {
		return fmt.Errorf("%s: %w", op, Errors.ErrTransactionNotFound)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	accountDto, err := s.repAccDto.CheckExistsAccount(ctx, transaction.AccountID)
	if errors.Is(err, Errors.ErrAccountNotFound) {
		return fmt.Errorf("%s: %w", op, Errors.ErrAccountNotFound)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	amount, err := common.ValidateCurrency(transaction.Currency, accountDto.Currency, transaction.Amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	switch {
	case transaction.ToAccount == 0:
		if err = s.repAccDto.UpdateBalance(ctx, transaction.AccountID, accountDto.Balance+amount); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if err = s.repTransaction.UpdateTransactionStatus(ctx, transaction.ID, response.StatusSuccess); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	case transaction.ToAccount > 0:
		if amount > accountDto.Balance {
			if err = s.repTransaction.UpdateTransactionStatus(ctx, transaction.ID, response.StatusError); err != nil {
				return fmt.Errorf("%s: %w", op, err)
			}
			return fmt.Errorf("%s: %w", op, Errors.ErrNegativeBalance)
		}
		if err = s.repAccDto.UpdateBalance(ctx, transaction.AccountID, accountDto.Balance-amount); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		if err = s.repTransaction.UpdateTransactionStatus(ctx, transaction.ID, response.StatusSuccess); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	default:
		if err = s.repTransaction.UpdateTransactionStatus(ctx, transaction.ID, response.StatusError); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		return fmt.Errorf("%s: %w", op, Errors.ErrIncorrectID)
	}

	return nil
}

func (s *Service) DeleteTransactionByID(ctx context.Context, id uint64) error {
	const op = "domain/transaction.Service.DeleteTransactionByID"

	err := s.repTransaction.DeleteTransactionByID(ctx, id)
	if errors.Is(err, Errors.ErrTransactionNotFound) {
		return fmt.Errorf("%s: %w", op, Errors.ErrTransactionNotFound)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
