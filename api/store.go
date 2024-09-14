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
	DeleteAccountById(Id string) error
	GetAllAccounts() ([]*Account, error)
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
func (s *PostgresStore) DeleteAccountById(Id string) error {
	query := `DELETE FROM Accounts WHERE id = $1`
	s.db.Exec(query, Id)
	return nil
}
func (s *PostgresStore) GetAllAccounts() ([]*Account, error) {
	query := `SELECT * FROM Accounts`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account:=new(Account)
		err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.AccountNumber, &account.Balance)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}
func (s *PostgresStore) GetAccountById(Id string) (*Account, error) {
	query := `SELECT id, firstname, lastname, account_number, balance FROM Accounts WHERE id = $1`
	row := s.db.QueryRow(query, Id)

	account := &Account{}
	err := row.Scan(&account.ID, &account.FirstName, &account.LastName, &account.AccountNumber, &account.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No account found with the given Id
		}
		return nil, err
	}

	return account, nil
}
