package service

import (
	"context"
	"errors"
	"fmt"
	"task/common"
	"task/internal/domain/Errors"
	"task/internal/domain/account/dto"
	"task/internal/domain/transaction/entity"
	"task/internal/domain/transaction/repository"
)

type Repository interface {
	CreateDepositTransaction(ctx context.Context, transaction *entity.Transaction) error
	CreateWithdrawTransaction(ctx context.Context, transaction *entity.Transaction) error
	GetTransactionByID(ctx context.Context, id uint64) (*entity.Transaction, error)
	UpdateTransactionStatus(ctx context.Context, id uint64, status string) error
	DeleteTransactionByID(ctx context.Context, id uint64) error
	UpdateBalance(ctx context.Context, id uint64, amount float64) error
	CheckExistsAccount(ctx context.Context, id uint64) (*dto.RegistrationCommand, error)
}

type Service struct {
	repository Repository
}

func NewService(di *common.DependencyContainer) *Service {
	return &Service{
		repository: repository.NewPostgresRepository(di.Pool),
	}
}

func (s *Service) CreateDepositTransaction(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error) {
	const op = "domain/transaction.Service.CreateDepositTransaction"

	_, err := s.repository.GetTransactionByID(ctx, transaction.ID)

	switch err {
	case nil:
		return nil, fmt.Errorf("%s: %w", op, Errors.ErrAccountExists)

	case Errors.ErrTransactionNotFound:
		err = s.repository.CreateDepositTransaction(ctx, transaction)
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

	_, err := s.repository.GetTransactionByID(ctx, transaction.ID)

	switch err {
	case nil:
		return nil, fmt.Errorf("%s: %w", op, Errors.ErrAccountExists)

	case Errors.ErrTransactionNotFound:
		err = s.repository.CreateWithdrawTransaction(ctx, transaction)
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

	transaction, err := s.repository.GetTransactionByID(ctx, id)
	if errors.Is(err, Errors.ErrTransactionNotFound) {
		return nil, fmt.Errorf("%s: %w", op, Errors.ErrTransactionNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return transaction, nil
}

// TODO make transaction status changing
func (s *Service) UpdateTransactionStatus(ctx context.Context, id uint64) error {
	const op = "domain/transaction.Service.UpdateWithSuccess"

	transaction, err := s.repository.GetTransactionByID(ctx, id)
	if errors.Is(err, Errors.ErrTransactionNotFound) {
		return fmt.Errorf("%s: %w", op, Errors.ErrTransactionNotFound)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	accountDto, err := s.repository.CheckExistsAccount(ctx, transaction.AccountID)
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
		if err = s.repository.UpdateBalance(ctx, transaction.AccountID, accountDto.Balance+amount); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	case transaction.ToAccount > 0:
		if amount > accountDto.Balance {
			return fmt.Errorf("%s: %w", op, Errors.ErrNegativeBalance)
		}
		if err = s.repository.UpdateBalance(ctx, transaction.AccountID, accountDto.Balance-amount); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	default:
		return fmt.Errorf("%s: %w", op, Errors.ErrIncorrectID)
	}

	return nil
}

//func (s *Service) UpdateWithError(ctx context.Context, id uint64) error {
//	const op = "domain/transaction.Service.UpdateWithError"
//
//	err := s.repository.UpdateWithError(ctx, id)
//	if errors.Is(err, Errors.ErrTransactionNotFound) {
//		return fmt.Errorf("%s: %w", op, Errors.ErrTransactionNotFound)
//	}
//	if err != nil {
//		return fmt.Errorf("%s: %w", op, err)
//	}
//	return nil
//}

func (s *Service) DeleteTransactionByID(ctx context.Context, id uint64) error {
	const op = "domain/transaction.Service.DeleteTransactionByID"

	err := s.repository.DeleteTransactionByID(ctx, id)
	if errors.Is(err, Errors.ErrTransactionNotFound) {
		return fmt.Errorf("%s: %w", op, Errors.ErrTransactionNotFound)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
