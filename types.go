package main

import "math/rand"

type Account struct {
	ID            int     `json:"id"`
	FirstName     string  `json:"first_name"`
	LastName      string  `json:"last_name"`
	AccountNumber int64   `json:"account_number"`
	Balance       float64 `json:"balance"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		ID:            rand.Intn(10000),
		FirstName:     firstName,
		LastName:      lastName,
		AccountNumber: int64(rand.Intn(100000)),
	}
}
