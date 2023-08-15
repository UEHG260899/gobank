package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddress string
}

func NewApiServer(listenAddress string) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccount))

	log.Println("JSON API server running on port: ", s.listenAddress)

	http.ListenAndServe(s.listenAddress, router)
}

func (s *APIServer) handleAccount(writter http.ResponseWriter, request *http.Request) error {
	if request.Method == "GET" {
		return s.handleGetAccount(writter, request)
	}

	if request.Method == "POST" {
		return s.handleCreateAccount(writter, request)
	}

	if request.Method == "DELETE" {
		return s.handleDeleteAccount(writter, request)
	}

	return fmt.Errorf("method not allowed %s", request.Method)
}

func (s *APIServer) handleGetAccount(writter http.ResponseWriter, request *http.Request) error {
	id := mux.Vars(request)["id"]
	fmt.Println(id)
	return WriteJSON(writter, http.StatusOK, &Account{})
}

func (s *APIServer) handleCreateAccount(writter http.ResponseWriter, request *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(writter http.ResponseWriter, request *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(writter http.ResponseWriter, request *http.Request) error {
	return nil
}

func WriteJSON(writter http.ResponseWriter, status int, v any) error {
	writter.Header().Add("Content-Type", "application/json")
	writter.WriteHeader(status)
	return json.NewEncoder(writter).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(writter http.ResponseWriter, request *http.Request) {
		if err := f(writter, request); err != nil {
			WriteJSON(writter, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
