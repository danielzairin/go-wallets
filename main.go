package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/danielzairin/go-wallets/wallets"
	_ "modernc.org/sqlite"
)

var DB_NAME string = "database.sqlite"

func main() {
	setup()

	db, err := sql.Open("sqlite", DB_NAME)
	if err != nil {
		panic(err)
	}

	wal := wallets.NewWallet(db)

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		wallets, err := wal.List()
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, http.StatusText(http.StatusInternalServerError))
			return
		}

		fmt.Fprintf(w, "wallets = %v\n", wallets)
	})

	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, http.StatusText(http.StatusBadRequest))
			return
		}

		wallet, err := wal.Create(name)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, http.StatusText(http.StatusInternalServerError))
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, wallet.Name)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func setup() {
	db, err := sql.Open("sqlite", DB_NAME)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`DROP TABLE IF EXISTS wallets`)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE wallets (name TEXT PRIMARY KEY, amount INTEGER)`)
	if err != nil {
		panic(err)
	}
}
