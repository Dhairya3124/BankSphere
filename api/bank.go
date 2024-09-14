package api

import (
	"encoding/json"
	"net/http"
)

type BankServer struct {
	http.Handler
	store Storage
}

func NewBankServer(store Storage) *BankServer {
	b := new(BankServer)
	router := http.NewServeMux()
	router.Handle("/account", http.HandlerFunc(b.handleAccount))
	router.Handle("/account/{id}", http.HandlerFunc(b.handleAccountById))
	router.Handle("/update", http.HandlerFunc(b.updateAccountHandler))
	router.Handle("/delete/", http.HandlerFunc(b.deleteAccountHandler))
	router.Handle("/transfer", http.HandlerFunc(b.transferBalanceHandler))

	b.Handler = router
	b.store = store
	return b

}

func (b *BankServer) handleAccount(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		b.createAccountHandler(w, r)
	case http.MethodGet:
		b.getAllAccountsHandler(w, r)

	}

}
func (b *BankServer) handleAccountById(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		b.getAccountByIdHandler(w, r)
	case http.MethodDelete:
		b.deleteAccountHandler(w, r)

	}
}
func (b *BankServer) createAccountHandler(w http.ResponseWriter, r *http.Request) error {

	accountRequest := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(accountRequest); err != nil {
		return err
	}
	account, err := NewAccount(accountRequest.FirstName, accountRequest.LastName)
	if err != nil {
		return err
	}
	b.store.CreateAccount(account)
	return WriteJSON(w, 200, account)
}

func (b *BankServer) getAllAccountsHandler(w http.ResponseWriter, r *http.Request) error {
	accounts, err := b.store.GetAllAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, 200, accounts)

}
func (b *BankServer) getAccountByIdHandler(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	account, err := b.store.GetAccountById(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, 200, account)

}
func (b *BankServer) deleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
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
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
