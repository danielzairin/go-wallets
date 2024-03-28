package main

import (
	"fmt"
	"net/http"
)

func reject(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
	fmt.Fprintln(w, http.StatusText(code))
}

func internalError(w http.ResponseWriter, err error) {
	fmt.Printf("Internal Server Error: %s\n", err)
	reject(w, http.StatusInternalServerError)
}
