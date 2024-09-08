package main

import (
	"net/http"

	"github.com/Dhairya3124/BankSphere/api"
)

func main() {
	server := api.NewBankServer()
	http.ListenAndServe(":5000", server)
}
