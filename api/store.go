package api

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}
type Storage interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) error
	DeleteAccount(*Account) error
	GetAccountById(Id string) (*Account, error)
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "golangProjects"
	dbname   = "postgres"
)

func NewStorage() (*PostgresStore, error) {
	postgresInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", postgresInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil

}
func (s *PostgresStore) CreateAccount(account *Account) error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS Accounts(
	id INTEGER PRIMARY KEY,
	firstname VARCHAR(50),
	lastname VARCHAR(50),
	account_number INTEGER,
	balance INTEGER
	)`)
	if err != nil {
		return err
	}
	query := `INSERT INTO Accounts(id,firstname,lastname,account_number,balance) VALUES($1, $2, $3, $4, $5)`
	_, queryError := s.db.Exec(query, account.ID, account.FirstName, account.LastName, account.AccountNumber, account.Balance)
	if queryError != nil {
		return queryError
	}

	return nil
}
func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}
func (s *PostgresStore) DeleteAccount(*Account) error {
	return nil
}
func (s *PostgresStore) GetAccountById(Id string) (*Account, error) {
	return nil, nil
}
