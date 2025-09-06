package account

import (
	"context"

	"github.com/shopspring/decimal"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateAccount(ctx context.Context, acc *Account) error {
	return s.repo.Create(ctx, acc)
}

func (s *Service) GetAccount(ctx context.Context, id int64) (*Account, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Deposit(ctx context.Context, id int64, amount decimal.Decimal) error {
	return s.repo.Deposit(ctx, id, amount)
}

func (s *Service) Withdraw(ctx context.Context, id int64, amount decimal.Decimal) error {
	return s.repo.Withdraw(ctx, id, amount)
}
