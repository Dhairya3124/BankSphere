package api

import (
	"fmt"
	"math/rand"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID            int64
	FirstName     string
	LastName      string
	AccountNumber int64
	Balance       int64
	EncryptedPassword string
}

func NewAccount(FirstName, LastName,Password string) (*Account, error) {
	if len(FirstName) == 0 && len(LastName) == 0 {
		return nil, fmt.Errorf("FirstName and LastName cannot be blank")
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Account{
		FirstName:     FirstName,
		LastName:      LastName,
		ID:            int64(rand.Intn(1000)),
		AccountNumber: int64(rand.Intn(10000)),
		EncryptedPassword: string(encryptedPassword),
	}, nil

}

type CreateAccountRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Password string `json:"password"`
}
type LogMessage struct {
	Message string `json:"msg"`
}
type UpdateBalanceRequest struct {
	AccountNumber int64 `json:"accountnumber"`
	Amount        int64 `json:"amount"`
}
type TransferBalanceRequest struct {
	SourceAccountNumber      int64 `json:"sourceAccount"`
	DestinationAccountNumber int64 `json:"destinationAccount"`
	Amount                   int64 `json:"amount"`
}
type LoginRequest struct {
	AccountNumber   int64  `json:"accountnumber"`
	Password string `json:"password"`
}
type LoginResponse struct {
	AccountNumber   int64  `json:"accountnumber"`
	Token string `json:"token"`
}
func (a *Account) ValidPassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw)) == nil
}