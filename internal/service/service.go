package service

import (
	"context"
	db "github.com/tredoc/go-balances/db/sqlc"
	"github.com/tredoc/go-balances/internal/store"
)

type Service struct {
	store *store.Store
}

func New(store *store.Store) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetAllBalances() ([]db.Balance, error) {
	return s.store.GetAllBalances(context.Background())
}

func (s *Service) GetAllCurrencies() ([]db.Currency, error) {
	return s.store.GetAllCurrencies(context.Background())
}

func (s *Service) GetAllEntries() ([]db.Entry, error) {
	return s.store.GetAllEntries(context.Background())
}

func (s *Service) GetAllTransfers() ([]db.Transfer, error) {
	return s.store.GetAllTransfers(context.Background())
}

func (s *Service) GetAllUsers() ([]db.User, error) {
	return s.store.GetAllUsers(context.Background())
}
