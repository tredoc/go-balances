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
