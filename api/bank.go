package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	router.Handle("/update", http.HandlerFunc(b.handleBalanceUpdate))
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
func (b *BankServer) handleBalanceUpdate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		b.updateAccountBalanceHandler(w, r)

	}
}
func (b *BankServer) createAccountHandler(w http.ResponseWriter, r *http.Request) error {

	accountRequest := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(accountRequest); err != nil {
		return WriteJSON(w, 400, logger(string(err.Error())))
	}
	account, err := NewAccount(accountRequest.FirstName, accountRequest.LastName)
	if err != nil {
		return WriteJSON(w, 400, logger(string(err.Error())))
	}
	b.store.CreateAccount(account)
	return WriteJSON(w, 200, account)
}

func (b *BankServer) getAllAccountsHandler(w http.ResponseWriter, r *http.Request) error {
	accounts, err := b.store.GetAllAccounts()
	if err != nil {
		return WriteJSON(w, 400, logger(string(err.Error())))
	}
	return WriteJSON(w, 200, accounts)

}
func (b *BankServer) getAccountByIdHandler(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return WriteJSON(w, 400, logger(string(err.Error())))
	}
	account, err := b.store.GetAccountById(id)
	if err != nil {
		return WriteJSON(w, 400, logger(string(err.Error())))
	}
	return WriteJSON(w, 200, account)

}
func (b *BankServer) deleteAccountHandler(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return WriteJSON(w, 400, logger(string(err.Error())))
	}
	err = b.store.DeleteAccountById(id)
	if err != nil {
		return WriteJSON(w, 400, logger(string(err.Error())))
	}
	return WriteJSON(w, 200, logger("Account Deleted"))

}
func (b *BankServer) updateAccountBalanceHandler(w http.ResponseWriter, r *http.Request) error {
	updateBalanceRequest := new(UpdateBalanceRequest)
	if err := json.NewDecoder(r.Body).Decode(updateBalanceRequest); err != nil {
		return WriteJSON(w, 400, logger(string(err.Error())))
	}
	defer r.Body.Close()
	err := b.store.UpdateAccountBalance(updateBalanceRequest.AccountNumber, updateBalanceRequest.Amount)
	if err != nil {
		return WriteJSON(w, 400, logger(string(err.Error())))
	}
	return WriteJSON(w, 200, logger("Account Updated"))
}
func (b *BankServer) transferBalanceHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("{}")
}
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
func getID(r *http.Request) (int, error) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}
func logger(msg string) *LogMessage {
	logs := new(LogMessage)
	logs.Message = msg
	return logs
}
