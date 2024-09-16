package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v5"
)

type BankServer struct {
	http.Handler
	store Storage
}

func NewBankServer(store Storage) *BankServer {
	b := new(BankServer)
	router := http.NewServeMux()
	router.Handle("/account", http.HandlerFunc(b.handleAccount))
	router.Handle("/account/{id}", withJWTAuth(http.HandlerFunc(b.handleAccountById),store))
	router.Handle("/update", http.HandlerFunc(b.handleBalanceUpdate))
	router.Handle("/transfer", http.HandlerFunc(b.handleBalanceTransfer))
	router.HandleFunc("/login", http.HandlerFunc(b.handleLoginAccount))

	b.Handler = router
	b.store = store
	return b

}

func (b *BankServer) handleLoginAccount(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "POST" {
		b.handleLogin(w,r)
	}
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
func (b *BankServer) handleBalanceTransfer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		b.transferBalanceHandler(w, r)
	}
}
func (b *BankServer) createAccountHandler(w http.ResponseWriter, r *http.Request) error {

	accountRequest := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(accountRequest); err != nil {
		return WriteJSON(w, 400, logger(string(err.Error())))
	}
	account, err := NewAccount(accountRequest.FirstName, accountRequest.LastName,accountRequest.Password)
	if err != nil {
		return WriteJSON(w, 400, logger(string(err.Error())))
	}
	err = b.store.CreateAccount(account)
	if err != nil {
		return WriteJSON(w,400,logger(string(err.Error())))
	}
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
func (b *BankServer) transferBalanceHandler(w http.ResponseWriter, r *http.Request) error {
	transferBalanceRequest := new(TransferBalanceRequest)
	if err := json.NewDecoder(r.Body).Decode(transferBalanceRequest); err != nil {
		return WriteJSON(w, 400, logger(string(err.Error())))
	}
	defer r.Body.Close()
	err := b.store.TransferBalancetoAccounts(transferBalanceRequest.SourceAccountNumber, transferBalanceRequest.DestinationAccountNumber, transferBalanceRequest.Amount)
	if err != nil {
		return WriteJSON(w, 400, logger(string(err.Error())))
	}
	return WriteJSON(w, 200, logger("Balance Transferred"))

}

func (b *BankServer) handleLogin(w http.ResponseWriter, r *http.Request) error {

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return WriteJSON(w,400,"invalid Body")
	}

	account, err := b.store.GetAccountByAccountNumber(int64(req.AccountNumber))
	if err != nil {
		return WriteJSON(w,400,"invalid AccountNumber")
	}

	if !account.ValidPassword(req.Password) {
		return WriteJSON(w,400,"invalid password")
	}

	token, err := createJWT(account)
	if err != nil {
		return WriteJSON(w,400,"wrong password")
	}

	resp := LoginResponse{
		Token:  token,
		AccountNumber: account.AccountNumber,
	}

	return WriteJSON(w, http.StatusOK, resp)
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
func createJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt":     15000,
		"accountNumber": account.AccountNumber,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden,  logger("permission denied"))
}


func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}
		userID, err := getID(r)
		if err != nil {
			permissionDenied(w)
			return
		}
		account, err := s.GetAccountById(userID)
		if err != nil {
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if account.AccountNumber != int64(claims["accountNumber"].(float64)) {
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}