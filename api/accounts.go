package api

import "math/rand"

type Account struct {
	ID            int64
	FirstName     string
	LastName      string
	AccountNumber int64
	Balance       int64
}

func NewAccount(FirstName, LastName string) *Account {
	return &Account{
		FirstName:     FirstName,
		LastName:      LastName,
		ID:            int64(rand.Intn(1000)),
		AccountNumber: int64(rand.Intn(10000)),
	}

}
