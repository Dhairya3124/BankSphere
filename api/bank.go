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
	b.Handler = router
	return b

}
func (b *BankServer) createAccountHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("{}")
}
