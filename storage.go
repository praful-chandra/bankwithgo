package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {

	connStr := "user=postgres dbname=pqgotest password=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	return s.CreateAccountTable()
}
func (s *PostgresStore) CreateAccountTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS account(
		id serial primary key,
		first_name varchar(50),
		last_name varchar(50),
		number serial,
		balance serial,
		created_at timestamp
	)
	`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
	INSERT INTO 
	account (first_name, last_name, number, balance, created_at)
	VALUES ($1,$2,$3,$4,$5)
`
	resp, err := s.db.Query(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt)
	if err != nil {
		return err
	}
	log.Printf("%+v\n", resp)

	return nil
}
func (s *PostgresStore) DeleteAccount(id int) error {
	return nil
}
func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}
func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM account WHERE id=$1`, id)
	if err != nil {
		return nil, err
	}

	account := new(Account)
	rows.Next()
	if err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	); err != nil {
		{
			return nil, err
		}
	}

	return account, nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM account`)
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}

	for rows.Next() {
		account := new(Account)
		if err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt,
		); err != nil {
			{
				return nil, err
			}
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}
