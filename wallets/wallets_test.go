package wallets_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/danielzairin/go-wallets/wallets"
	_ "modernc.org/sqlite"
)

var db *sql.DB

func TestMain(m *testing.M) {
	setup()
	defer db.Close()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	var err error
	db, err = sql.Open("sqlite", "database_test.sqlite")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`DROP TABLE IF EXISTS wallets`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE wallets (name TEXT PRIMARY KEY, amount INTEGER)`)
	if err != nil {
		panic(err)
	}
}

func TestTopUp(t *testing.T) {
	walletsApp := wallets.NewWallet(db)

	wallet, err := walletsApp.Create("test")
	if err != nil {
		t.Error(err)
	}

	_, err = walletsApp.TopUp(wallet, 10)
	if err != nil {
		t.Error(err)
	}

	if wallet.Amount != 10 {
		t.Errorf("expected = %d, got = %d", 10, wallet.Amount)
	}
}
