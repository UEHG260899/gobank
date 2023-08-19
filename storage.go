package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountById(int) (*Account, error)
	GetAccountByNumber(int64) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	dbPassword := os.Getenv("DB_PASSWORD")
	connectionString := fmt.Sprintf("user=postgres dbname=postgres password=%s sslmode=disable", dbPassword)
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	statement := `CREATE TABLE IF NOT EXISTS accounts (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(255),
		last_name VARCHAR(255),
		encrypted_password VARCHAR,
		number INT,
		balance NUMERIC(2)
	);`

	_, err := s.db.Exec(statement)
	return err
}

func (s *PostgresStore) CreateAccount(account *Account) error {
	statement := `
		INSERT INTO accounts (first_name, last_name, number, encrypted_password, balance)
		VALUES ($1, $2, $3, $4, $5);
	`

	_, err := s.db.Exec(
		statement,
		account.FirstName,
		account.LastName,
		account.AccountNumber,
		account.EncryptedPassword,
		account.Balance,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteAccount(id int) error {
	statement := "DELETE FROM accounts WHERE id = $1;"
	_, err := s.db.Exec(statement, id)
	return err
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	query := `
		SELECT id,
			   first_name,
			   last_name,
			   number,
			   balance
		FROM accounts
		WHERE id = $1;
	`
	rows, err := s.db.Query(query, id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", id)
}

func (s *PostgresStore) GetAccountByNumber(number int64) (*Account, error) {
	query := `
		SELECT id,
			   first_name,
			   last_name,
			   number,
			   balance
		FROM accounts
		WHERE number = $1;
	`

	rows, err := s.db.Query(query, number)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf("account %d not found", number)
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := `
		SELECT id,
			   first_name,
			   last_name,
			   number,
			   balance
		FROM accounts;
	`

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		account, err := scanIntoAccount(rows)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(&account.ID, &account.FirstName, &account.LastName, &account.AccountNumber, &account.Balance)

	return account, err
}
