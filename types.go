package main

import (
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	AccountNumber int64  `json:"account_number"`
	Password      string `json:"password"`
}

type LoginReponse struct {
	Number int64  `json:"account_number"`
	Token  string `json:"token"`
}

type TransferRequest struct {
	ToAccount int64   `json:"to_account"`
	Amount    float64 `json:"amount"`
}

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type Account struct {
	ID                int     `json:"id"`
	FirstName         string  `json:"first_name"`
	LastName          string  `json:"last_name"`
	AccountNumber     int64   `json:"account_number"`
	EncryptedPassword string  `json:"-"`
	Balance           float64 `json:"balance"`
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		AccountNumber:     int64(rand.Intn(100000)),
		EncryptedPassword: string(encrypted),
	}, nil
}
