package service

import (
	"context"
	"errors"
	db "github.com/tredoc/go-balances/db/sqlc"
	"github.com/tredoc/go-balances/internal/store"
	"time"
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

func (s *Service) GetBalanceById(id uint64) (db.Balance, error) {
	return s.store.GetBalanceByID(context.Background(), id)
}

func (s *Service) Deposit(id uint64, amount int64) (*db.Balance, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	tx, err := s.store.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := s.store.WithTx(tx)

	balance, err := qtx.GetBalanceByIDForUpdate(context.Background(), id)
	if err != nil {
		return nil, err
	}

	_, err = qtx.CreateEntry(context.Background(), db.CreateEntryParams{BalanceID: id, Amount: amount})

	err = qtx.UpdateBalance(context.Background(), db.UpdateBalanceParams{ID: id, Amount: balance.Amount + amount})
	if err != nil {
		return nil, err
	}
	err = tx.Commit()

	balance.Amount += amount
	return &balance, err
}

func (s *Service) Withdraw(id uint64, amount int64) (*db.Balance, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	tx, err := s.store.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := s.store.WithTx(tx)

	balance, err := qtx.GetBalanceByIDForUpdate(context.Background(), id)
	if err != nil {
		return nil, err
	}

	if balance.Amount < amount {
		return nil, errors.New("insufficient funds")
	}

	_, err = qtx.CreateEntry(context.Background(), db.CreateEntryParams{BalanceID: id, Amount: -amount})
	if err != nil {
		return nil, err
	}

	err = qtx.UpdateBalance(context.Background(), db.UpdateBalanceParams{ID: id, Amount: balance.Amount - amount})
	if err != nil {
		return nil, err
	}
	err = tx.Commit()

	balance.Amount -= amount
	return &balance, err
}

func (s *Service) Transfer(fromID uint64, toID uint64, amount int64) (*db.Balance, *db.Balance, error) {
	if amount <= 0 {
		return nil, nil, errors.New("amount must be positive")
	}

	tx, err := s.store.DB.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	qtx := s.store.WithTx(tx)

	// always use same update order to avoid deadlock
	if fromID < toID {
		// using regular GetBalanceByID will cause deadlock
		balanceFrom, err := qtx.GetBalanceByIDForUpdate(context.Background(), fromID)
		if err != nil {
			return nil, nil, err
		}

		if balanceFrom.Amount < amount {
			return nil, nil, errors.New("insufficient funds")
		}

		// using regular GetBalanceByID will cause deadlock
		balanceTo, err := qtx.GetBalanceByIDForUpdate(context.Background(), toID)
		if err != nil {
			return nil, nil, err
		}

		_, err = qtx.CreateTransfer(context.Background(), db.CreateTransferParams{FromBalanceID: fromID, ToBalanceID: toID, Amount: amount})
		if err != nil {
			return nil, nil, err
		}

		// add sleep to emulate slow db
		time.Sleep(150 * time.Millisecond)
		err = qtx.UpdateBalance(context.Background(), db.UpdateBalanceParams{ID: fromID, Amount: balanceFrom.Amount - amount})
		if err != nil {
			return nil, nil, err
		}

		// add sleep to emulate slow db
		time.Sleep(100 * time.Millisecond)
		err = qtx.UpdateBalance(context.Background(), db.UpdateBalanceParams{ID: toID, Amount: balanceTo.Amount + amount})
		if err != nil {
			return nil, nil, err
		}

		err = tx.Commit()

		balanceFrom.Amount -= amount
		balanceTo.Amount += amount
		return &balanceFrom, &balanceTo, err

	} else {

		// using regular GetBalanceByID will cause deadlock
		balanceTo, err := qtx.GetBalanceByIDForUpdate(context.Background(), toID)
		if err != nil {
			return nil, nil, err
		}

		// using regular GetBalanceByID will cause deadlock
		balanceFrom, err := qtx.GetBalanceByIDForUpdate(context.Background(), fromID)
		if err != nil {
			return nil, nil, err
		}

		if balanceFrom.Amount < amount {
			return nil, nil, errors.New("insufficient funds")
		}

		_, err = qtx.CreateTransfer(context.Background(), db.CreateTransferParams{FromBalanceID: fromID, ToBalanceID: toID, Amount: amount})
		if err != nil {
			return nil, nil, err
		}

		// add sleep to emulate slow db
		time.Sleep(150 * time.Millisecond)
		err = qtx.UpdateBalance(context.Background(), db.UpdateBalanceParams{ID: toID, Amount: balanceTo.Amount + amount})
		if err != nil {
			return nil, nil, err
		}

		// add sleep to emulate slow db
		time.Sleep(100 * time.Millisecond)
		err = qtx.UpdateBalance(context.Background(), db.UpdateBalanceParams{ID: fromID, Amount: balanceFrom.Amount - amount})
		if err != nil {
			return nil, nil, err
		}

		err = tx.Commit()

		balanceFrom.Amount -= amount
		balanceTo.Amount += amount
		return &balanceFrom, &balanceTo, err
	}
}

func (s *Service) GetLastTransferID() (uint64, error) {
	return s.store.GetLastTransferID(context.Background())
}

func (s *Service) GetAllCurrencies() ([]db.Currency, error) {
	return s.store.GetAllCurrencies(context.Background())
}

func (s *Service) GetAllEntries() ([]db.Entry, error) {
	return s.store.GetAllEntries(context.Background())
}

func (s *Service) GetLastEntryID() (uint64, error) {
	return s.store.GetLastEntryID(context.Background())
}

func (s *Service) GetAllTransfers() ([]db.Transfer, error) {
	return s.store.GetAllTransfers(context.Background())
}

func (s *Service) GetAllUsers() ([]db.User, error) {
	return s.store.GetAllUsers(context.Background())
}
