package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/tredoc/go-balances/internal/service"
	"github.com/tredoc/go-balances/internal/store"
	"log"
	"os"
	"sync"
	"testing"
)

var conn *sql.DB
var services *service.Service
var storage *store.Store

func TestMain(m *testing.M) {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	if user == "" || pass == "" || host == "" || port == "" || name == "" {
		log.Fatal("missing environment variables")
	}

	var err error
	conn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, name))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	err = conn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	storage = store.New(conn)
	services = service.New(storage)

	m.Run()
}

func TestStore(t *testing.T) {
	t.Run("Test store creation", func(t *testing.T) {
		assert.NotEqual(t, storage, &store.Store{})
	})
}

func TestService(t *testing.T) {
	t.Run("Test service creation", func(t *testing.T) {
		assert.NotEqual(t, services, &service.Service{})
	})
}

func TestDeposit(t *testing.T) {
	t.Run("Test single deposit", func(t *testing.T) {
		balanceID := uint64(1)
		balance, err := services.GetBalanceById(balanceID)
		assert.Nil(t, err)

		amount := int64(100)
		balanceUPD, err := services.Deposit(balanceID, amount)
		assert.Nil(t, err)
		assert.Equal(t, balanceUPD.Amount, balance.Amount+amount)
	})

	t.Run("Test concurrent deposit", func(t *testing.T) {
		balanceID := uint64(4)
		balance, err := services.GetBalanceById(balanceID)
		assert.Nil(t, err)

		lastEntryID, err := services.GetLastEntryID()
		assert.Nil(t, err)

		amount := int64(5)
		times := 10

		var wg sync.WaitGroup
		for i := 0; i < times; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := services.Deposit(balanceID, amount)
				assert.Nil(t, err)
			}()
		}
		wg.Wait()

		lastEntryIDNew, err := services.GetLastEntryID()
		assert.Nil(t, err)
		assert.Equal(t, lastEntryIDNew, lastEntryID+uint64(times))

		balanceUPD, err := services.GetBalanceById(balanceID)
		assert.Nil(t, err)
		assert.Equal(t, balanceUPD.Amount, balance.Amount+amount*int64(times))
	})
}

func TestWithdraw(t *testing.T) {
	t.Run("Test single successful withdraw", func(t *testing.T) {
		balanceID := uint64(1)
		balance, err := services.GetBalanceById(balanceID)
		assert.Nil(t, err)

		amount := int64(100)
		balanceUPD, err := services.Withdraw(balanceID, amount)
		assert.Nil(t, err)
		assert.Equal(t, balanceUPD.Amount, balance.Amount-amount)

		amount = int64(1<<63 - 1)
		_, err = services.Withdraw(1, amount)
		assert.Error(t, err)
	})

	t.Run("Test single overbalance withdraw", func(t *testing.T) {
		balanceID := uint64(1)
		amount := int64(1<<63 - 1)
		_, err := services.Withdraw(balanceID, amount)
		assert.Error(t, err)
	})

	t.Run("Test concurrent withdraw", func(t *testing.T) {
		balanceID := uint64(6)
		balance, err := services.GetBalanceById(balanceID)
		assert.Nil(t, err)

		lastEntryID, err := services.GetLastEntryID()
		assert.Nil(t, err)

		amount := int64(5)
		times := 10

		var wg sync.WaitGroup
		for i := 0; i < times; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := services.Withdraw(balanceID, amount)
				assert.Nil(t, err)
			}()
		}
		wg.Wait()

		lastEntryIDNew, err := services.GetLastEntryID()
		assert.Nil(t, err)
		assert.Equal(t, lastEntryIDNew, lastEntryID+uint64(times))

		balanceUPD, err := services.GetBalanceById(balanceID)
		assert.Nil(t, err)
		assert.Equal(t, balanceUPD.Amount, balance.Amount-amount*int64(times))
	})
}

func TestTransfer(t *testing.T) {
	t.Run("Test single transfer", func(t *testing.T) {
		fromID := uint64(2)
		toID := uint64(10)

		balanceFrom, err := services.GetBalanceById(fromID)
		assert.Nil(t, err)

		balanceTo, err := services.GetBalanceById(toID)
		assert.Nil(t, err)

		amount := int64(10)
		balanceFromUPD, balanceToUPD, err := services.Transfer(fromID, toID, amount)
		assert.Nil(t, err)
		assert.Equal(t, balanceFromUPD.Amount, balanceFrom.Amount-amount)
		assert.Equal(t, balanceToUPD.Amount, balanceTo.Amount+amount)

		amount = int64(1<<63 - 1)
		_, _, err = services.Transfer(fromID, toID, amount)
		assert.Error(t, err)
	})
}
