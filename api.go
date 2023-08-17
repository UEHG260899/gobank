package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddress string
	store         Storage
}

func NewApiServer(listenAddress string, store Storage) *APIServer {
	return &APIServer{
		listenAddress: listenAddress,
		store:         store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleAccountWithId), s.store))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer))

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

	return fmt.Errorf("method not allowed %s", request.Method)
}

func (s *APIServer) handleAccountWithId(writter http.ResponseWriter, request *http.Request) error {
	if request.Method == "GET" {
		return s.handleGetAccountById(writter, request)
	}

	if request.Method == "DELETE" {
		return s.handleDeleteAccount(writter, request)
	}

	return fmt.Errorf("method not allowed %s", request.Method)
}

func (s *APIServer) handleGetAccount(writter http.ResponseWriter, request *http.Request) error {
	result, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJSON(writter, http.StatusOK, result)
}

func (s *APIServer) handleGetAccountById(writter http.ResponseWriter, request *http.Request) error {
	idAsInt, err := getIdFromParameters(request)

	if err != nil {
		return err
	}

	account, err := s.store.GetAccountById(idAsInt)

	if err != nil {
		return err
	}

	return WriteJSON(writter, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(writter http.ResponseWriter, request *http.Request) error {
	accountRequest := new(CreateAccountRequest)

	if err := json.NewDecoder(request.Body).Decode(accountRequest); err != nil {
		return err
	}

	account := NewAccount(accountRequest.FirstName, accountRequest.LastName)

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	tokenString, err := createJWT(account)

	if err != nil {
		return err
	}

	fmt.Println(tokenString)

	return WriteJSON(writter, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(writter http.ResponseWriter, request *http.Request) error {
	idAsInt, err := getIdFromParameters(request)

	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(idAsInt); err != nil {
		return err
	}

	return WriteJSON(writter, http.StatusOK, map[string]int{"deleted": idAsInt})
}

func (s *APIServer) handleTransfer(writter http.ResponseWriter, request *http.Request) error {
	transferRequest := new(TransferRequest)

	if err := json.NewDecoder(request.Body).Decode(transferRequest); err != nil {
		return err
	}

	defer request.Body.Close()
	return WriteJSON(writter, http.StatusOK, transferRequest)
}

func WriteJSON(writter http.ResponseWriter, status int, v any) error {
	writter.Header().Add("Content-Type", "application/json")
	writter.WriteHeader(status)
	return json.NewEncoder(writter).Encode(v)
}

func createJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt":     15000,
		"accountNumber": account.AccountNumber,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func withJWTAuth(handlerFunc http.HandlerFunc, storage Storage) http.HandlerFunc {
	return func(writter http.ResponseWriter, request *http.Request) {
		tokenString := request.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)
		claims := token.Claims.(jwt.MapClaims)

		if err != nil {
			WriteJSON(writter, http.StatusForbidden, ApiError{Error: "invalid token"})
			return
		}

		if !token.Valid {
			WriteJSON(writter, http.StatusForbidden, ApiError{Error: "invalid token"})
			return
		}

		userId, err := getIdFromParameters(request)

		if err != nil {
			WriteJSON(writter, http.StatusForbidden, ApiError{Error: "invalid token"})
			return
		}

		account, err := storage.GetAccountById(userId)

		if account.AccountNumber != int64(claims["accountNumber"].(float64)) {
			WriteJSON(writter, http.StatusForbidden, ApiError{"Access denied"})
			return
		}

		if err != nil {
			WriteJSON(writter, http.StatusForbidden, ApiError{Error: "invalid token"})
			return
		}

		handlerFunc(writter, request)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(writter http.ResponseWriter, request *http.Request) {
		if err := f(writter, request); err != nil {
			WriteJSON(writter, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getIdFromParameters(request *http.Request) (int, error) {
	idAsString := mux.Vars(request)["id"]
	idAsInt, err := strconv.Atoi(idAsString)

	if err != nil {
		return idAsInt, fmt.Errorf("Invalid id given %s", idAsString)
	}

	return idAsInt, nil
}
