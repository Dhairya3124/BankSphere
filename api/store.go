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
func (s *PostgresStore) CreateAccount(*Account) error {
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
