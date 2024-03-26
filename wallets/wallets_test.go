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
	w := wallets.NewWallet(db)

	wallet, err := w.Create("test")
	if err != nil {
		t.Error(err)
	}

	_, err = w.TopUp(wallet, 10)
	if err != nil {
		t.Error(err)
	}

	if wallet.Amount != 10 {
		t.Errorf("expected = %d, got = %d", 10, wallet.Amount)
	}
}

func TestList(t *testing.T) {
	w := wallets.NewWallet(db)

	walletList, err := w.List()
	if err != nil {
		t.Error(err)
	}

	if len(walletList) != 1 {
		t.Errorf("expected len(walletList) = %d, got = %d", 1, len(walletList))
	}
}

func TestFindByID(t *testing.T) {
	w := wallets.NewWallet(db)

	wallet, err := w.FindByName("test")
	if err != nil {
		t.Error(err)
	}

	if wallet.Name != "test" {
		t.Errorf("expected wallet.Name = %s, got = %s", "test", wallet.Name)
	}
}

func TestTransfer(t *testing.T) {
	w := wallets.NewWallet(db)

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
