package main

import (
	"math/rand"
)

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Account struct {
	ID            int     `json:"id"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	AccountNumber int64   `json:"account_number"`
	Balance       float64 `json:"balance"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		FirstName:     firstName,
		LastName:      lastName,
		AccountNumber: int64(rand.Intn(100000)),
	}
}

// docker run --name some-postgres -e POSTGRES_PASSWORD=He26Gu0899 -p 5432:5432 -d postgres
