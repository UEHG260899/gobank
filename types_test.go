package main

import (
	"fmt"
	"testing"
)

func TestNewAccount(t *testing.T) {
	account, err := NewAccount("a", "b", "hunter")

	if err != nil {
		t.Error("There shouldnt be an error")
	}

	fmt.Printf("%+v\n", account)
}
