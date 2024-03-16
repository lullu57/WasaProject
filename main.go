package main

import (
	"WASAPROJECT/services"
	"net/http"
)

func main() {
	srv := services.NewServer()
	http.ListenAndServe(":8080", srv)
}
