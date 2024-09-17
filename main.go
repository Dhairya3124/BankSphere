package main

import (
	"log"
	"net/http"

	"github.com/Dhairya3124/BankSphere/api"
)

func main() {
	store, err := api.NewStorage()
	if err != nil {
		log.Fatal(err)
	}
	server := api.NewBankServer(store)
	log.Fatal(http.ListenAndServe(":5000", server))
	
}
