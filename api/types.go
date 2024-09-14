package api

import (
	"fmt"
	"math/rand"
)

type Account struct {
	ID            int64
	FirstName     string
	LastName      string
	AccountNumber int64
	Balance       int64
}

func NewAccount(FirstName, LastName string) (*Account, error) {
	if len(FirstName) == 0 && len(LastName) == 0 {
		return nil, fmt.Errorf("FirstName and LastName cannot be blank")
	}
	return &Account{
		FirstName:     FirstName,
		LastName:      LastName,
		ID:            int64(rand.Intn(1000)),
		AccountNumber: int64(rand.Intn(10000)),
	}, nil

}

type CreateAccountRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}
type LogMessage struct {
	Message string `json:"msg"`
}
