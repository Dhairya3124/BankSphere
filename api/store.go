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
	UpdateAccountBalance(accountNumber, amount int64) error
	DeleteAccountById(Id int) error
	GetAllAccounts() ([]*Account, error)
	GetAccountById(Id int) (*Account, error)
	TransferBalancetoAccounts(sourceAccountNumber, destinationAccountNumber, amount int64) error
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
func (s *PostgresStore) UpdateAccountBalance(accountNumber, amount int64) error {
	query := `SELECT id FROM Accounts WHERE account_number = $1`
	row, rowError := s.db.Query(query, accountNumber)
	if !row.Next() || rowError != nil {
		return fmt.Errorf("account not found with account number as %d", accountNumber)
	}
	query = `UPDATE Accounts SET balance = balance + $2 WHERE account_number = $1`
	_, err := s.db.Exec(query, accountNumber, amount)
	if err != nil {
		return err
	}
	return nil
}
func (s *PostgresStore) DeleteAccountById(Id int) error {
	query := `SELECT id FROM Accounts WHERE id = $1`
	row, rowError := s.db.Query(query, Id)
	if !row.Next() || rowError != nil {
		return fmt.Errorf("account not found with id as %d", Id)
	}
	query = `DELETE FROM Accounts WHERE id = $1`
	_, err := s.db.Exec(query, Id)
	if err != nil {
		return fmt.Errorf("account not found with id as %d", Id)
	}
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
		account := new(Account)
		err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.AccountNumber, &account.Balance)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}
func (s *PostgresStore) GetAccountById(Id int) (*Account, error) {
	query := `SELECT id, firstname, lastname, account_number, balance FROM Accounts WHERE id = $1`
	row := s.db.QueryRow(query, Id)

	account := &Account{}
	err := row.Scan(&account.ID, &account.FirstName, &account.LastName, &account.AccountNumber, &account.Balance)
	if err != nil {
		return nil, fmt.Errorf("account not found with id as %d", Id) // No account found with the given Id

	}

	return account, nil
}
func (s *PostgresStore) TransferBalancetoAccounts(sourceAccountNumber, destinationAccountNumber int64, amount int64) error {
	query := `SELECT id FROM Accounts WHERE account_number = $1`
	row := s.db.QueryRow(query, destinationAccountNumber)
	var destAccountId int
	err := row.Scan(&destAccountId)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("destination account not found with account number %d", destinationAccountNumber)
		}
		return fmt.Errorf("error querying destination account: %v", err)
	}

	query = `SELECT balance FROM Accounts WHERE account_number = $1`
	row = s.db.QueryRow(query, sourceAccountNumber)
	var balance int64
	err = row.Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("source account not found with account number %d", sourceAccountNumber)
		}
		return fmt.Errorf("error querying source account balance: %v", err)
	}

	if balance < amount {
		return fmt.Errorf("insufficient balance: current balance is %d, but %d is required", balance, amount)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}

	_, err = tx.Exec(`UPDATE Accounts SET balance = balance - $1 WHERE account_number = $2`, amount, sourceAccountNumber)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error updating source account balance: %v", err)
	}

	_, err = tx.Exec(`UPDATE Accounts SET balance = balance + $1 WHERE account_number = $2`, amount, destinationAccountNumber)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error updating destination account balance: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}
