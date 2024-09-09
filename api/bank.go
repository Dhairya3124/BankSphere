package api

import (
	"encoding/json"
	"net/http"
)

type BankServer struct {
	http.Handler
}

func NewBankServer() *BankServer {
	b := new(BankServer)
	router := http.NewServeMux()
	router.Handle("/create", http.HandlerFunc(b.createAccountHandler))
	router.Handle("/get", http.HandlerFunc(b.getAccountHandler))
	router.Handle("/update", http.HandlerFunc(b.updateAccountHandler))
	router.Handle("/delete", http.HandlerFunc(b.deleteAccountHandler))
	router.Handle("/transfer", http.HandlerFunc(b.transferBalanceHandler))

	b.Handler = router
	return b

}
func (b *BankServer) createAccountHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("{}")
}
func (b *BankServer) getAccountHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("{}")
}
func (b *BankServer) deleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("{}")
}
func (b *BankServer) updateAccountHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("{}")
}
func (b *BankServer) transferBalanceHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("{}")
}