package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/danielzairin/go-wallets/wallets"
	_ "modernc.org/sqlite"
)

var PORT int = 8080
var DB_NAME string = "database.sqlite"

type App struct {
	wallets *wallets.Wallets
}

func main() {
	db, err := sql.Open("sqlite", DB_NAME)
	if err != nil {
		panic(err)
	}

	app := App{
		wallets: wallets.NewWallet(db),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", app.listWallets)
	mux.HandleFunc("POST /", app.postCreateWallet)
	mux.HandleFunc("POST /top-up/{walletName}", app.postTopUpWallet)
	mux.HandleFunc("POST /transfer", app.postTransferFunds)

	log.Printf("Listening on port %d", PORT)
	err = http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux)
	if err != nil {
		log.Fatal(err)
	}
}
