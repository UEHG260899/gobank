package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func seedAccount(store Storage, firstName, lastName, password string) *Account {
	acc, err := NewAccount(firstName, lastName, password)

	if err != nil {
		log.Fatal(err)
	}

	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	return acc
}

func seedAccounts(store Storage) {
	seedAccount(store, "Uriel", "Hernandez", "pass23456")
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Could not load Env File")
	}

	seed := flag.Bool("seed", false, "seed the db")
	flag.Parse()

	store, err := NewPostgresStore()

	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("seeding the database")
		seedAccounts(store)
	}

	server := NewApiServer(":3000", store)
	server.Run()
}
