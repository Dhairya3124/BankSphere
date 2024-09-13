package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

type BankServer struct {
	http.Handler
	store Storage
}

func NewBankServer(store Storage) *BankServer {
	b := new(BankServer)
	router := http.NewServeMux()
	router.Handle("/create", http.HandlerFunc(b.createAccountHandler))
	router.Handle("/get/", http.HandlerFunc(b.getAccountHandler))
	router.Handle("/update", http.HandlerFunc(b.updateAccountHandler))
	router.Handle("/delete/", http.HandlerFunc(b.deleteAccountHandler))
	router.Handle("/transfer", http.HandlerFunc(b.transferBalanceHandler))

	b.Handler = router
	b.store = store
	return b

}

const jsonContentType = "application/json"

func (b *BankServer) createAccountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	a := NewAccount("Jack", "Black")
	b.store.CreateAccount(a)

	json.NewEncoder(w).Encode(a)
}
func (b *BankServer) getAccountHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/get/")
	a, err := b.store.GetAccountById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(a)
	}

}
func (b *BankServer) deleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/delete/")
	err := b.store.DeleteAccountById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(err)
	} 
}
func (b *BankServer) updateAccountHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("{}")
}
func (b *BankServer) transferBalanceHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("{}")
}
