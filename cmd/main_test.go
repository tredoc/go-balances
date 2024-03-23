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
	assert.NotEqual(t, storage, &store.Store{})
}

func TestService(t *testing.T) {
	assert.NotEqual(t, services, &service.Service{})
}

func TestDeposit(t *testing.T) {
	balance, err := services.GetBalanceById(1)
	assert.Nil(t, err)

	deposit := int64(100)
	balanceUPD, err := services.Deposit(1, deposit)
	assert.Nil(t, err)
	assert.Equal(t, balanceUPD.Amount, balance.Amount+deposit)
}

func TestWithdraw(t *testing.T) {
	balance, err := services.GetBalanceById(1)
	assert.Nil(t, err)

	withdraw := int64(100)
	balanceUPD, err := services.Withdraw(1, withdraw)
	assert.Nil(t, err)
	assert.Equal(t, balanceUPD.Amount, balance.Amount-withdraw)

	withdraw = int64(1<<63 - 1)
	_, err = services.Withdraw(1, withdraw)
	assert.Error(t, err)
}

func TestTransfer(t *testing.T) {
	fromID := uint64(2)
	toID := uint64(10)

	balanceFrom, err := services.GetBalanceById(fromID)
	assert.Nil(t, err)

	balanceTo, err := services.GetBalanceById(toID)
	assert.Nil(t, err)

	amount := int64(100)
	balanceFromUPD, balanceToUPD, err := services.Transfer(fromID, toID, amount)
	assert.Nil(t, err)
	assert.Equal(t, balanceFromUPD.Amount, balanceFrom.Amount-amount)
	assert.Equal(t, balanceToUPD.Amount, balanceTo.Amount+amount)

	amount = int64(1<<63 - 1)
	_, _, err = services.Transfer(fromID, toID, amount)
	assert.Error(t, err)
}
