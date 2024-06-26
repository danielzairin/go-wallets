package wallets_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/danielzairin/go-wallets/internal/wallets"
	_ "modernc.org/sqlite"
)

var db *sql.DB
var w *wallets.Wallets

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

	w = wallets.NewWallet(db)
}

func TestTopUp(t *testing.T) {
	wallet, err := w.Create("test")
	if err != nil {
		t.Error(err)
	}

	_, err = w.TopUp(wallet, 30)
	if err != nil {
		t.Error(err)
	}

	if wallet.Amount != 30 {
		t.Errorf("expected = %d, got = %d", 30, wallet.Amount)
	}
}

func TestList(t *testing.T) {
	walletList, err := w.List()
	if err != nil {
		t.Error(err)
	}

	if len(walletList) != 1 {
		t.Errorf("expected len(walletList) = %d, got = %d", 1, len(walletList))
	}
}

func TestFindByID(t *testing.T) {
	wallet, err := w.FindByName("test")
	if err != nil {
		t.Error(err)
	}

	if wallet.Name != "test" {
		t.Errorf("expected wallet.Name = %s, got = %s", "test", wallet.Name)
	}
}

func TestTransfer(t *testing.T) {
	wallet1, err := w.FindByName("test")
	if err != nil {
		t.Error(err)
	}

	wallet2, err := w.Create("dummy")
	if err != nil {
		t.Error(err)
	}

	wallet1Before := wallet1.Amount
	wallet2Before := wallet2.Amount
	transferAmount := 5
	w.Transfer(&wallet1, wallet2, transferAmount)

	if wallet2.Amount != (wallet2Before + transferAmount) {
		t.Errorf("wallet2.Amount received incorrect funds, got = %d, want  = %d", wallet2.Amount, wallet2Before+transferAmount)
	}

	if wallet1.Amount != (wallet1Before - transferAmount) {
		t.Errorf("wallet1.Amount sent incorrect funds, got = %d, want  = %d", wallet1.Amount, wallet1Before-transferAmount)
	}
}
