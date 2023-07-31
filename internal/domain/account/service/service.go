package service

import (
	"context"
	"errors"
	"fmt"
	"task/common"
	"task/internal/domain/Errors"
	"task/internal/domain/account/entity"
	"task/internal/domain/account/repository"
)

type Repository interface {
	Save(ctx context.Context, account *entity.Account) error
	Get(ctx context.Context, id uint64) (*entity.Account, error)
	Delete(ctx context.Context, id uint64) error
	Update(ctx context.Context, id uint64, balance float64, currency string) error
}

type Service struct {
	repository Repository
}

func NewService(di *common.DependencyContainer) *Service {
	return &Service{
		repository: repository.NewPostgresRepository(di.Pool),
	}
}

func (s *Service) SaveAccount(ctx context.Context, account *entity.Account) (*entity.Account, error) {
	const op = "domain/account.Service.SaveAccount"

	_, err := s.repository.Get(ctx, account.ID)

	switch err {
	case nil:
		return nil, fmt.Errorf("%s: %w", op, Errors.ErrAccountExists)

	case Errors.ErrAccountNotFound:
		err = s.repository.Save(ctx, account)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		return account, nil
	default:
		return nil, fmt.Errorf("%s: %w", op, err)
	}

}

func (s *Service) GetAccount(ctx context.Context, id uint64) (*entity.Account, error) {
	const op = "domain/account.Service.GetAccount"

	account, err := s.repository.Get(ctx, id)
	if errors.Is(err, Errors.ErrAccountNotFound) {
		return nil, fmt.Errorf("%s: %w", op, Errors.ErrAccountNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return account, nil
}

func (s *Service) DeleteAccount(ctx context.Context, id uint64) error {
	const op = "domain/account.Service.Delete"

	err := s.repository.Delete(ctx, id)
	if errors.Is(err, Errors.ErrAccountNotFound) {
		return fmt.Errorf("%s: %w", op, Errors.ErrAccountNotFound)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Service) UpdateBalance(ctx context.Context, id uint64, balance float64, currency string) error {
	const op = "domain/account.Service.Update"

	err := s.repository.Update(ctx, id, balance, currency)
	if errors.Is(err, Errors.ErrAccountNotFound) {
		return fmt.Errorf("%s: %w", op, Errors.ErrAccountNotFound)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
