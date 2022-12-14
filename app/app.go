package app

import (
	"context"
	"fmt"
	"github.com/akshanshgusain/Hexagonal-Architecture/domain"
	"github.com/akshanshgusain/Hexagonal-Architecture/service"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os"
)

// Sanity test is a test of newly deployed environment.
// responsible to check if all the env variables are passed correctly
func sanityCheck() {
	if os.Getenv("SERVER_ADDRESS") == "" ||
		os.Getenv("SERVER_PORT") == "" {
		log.Fatal("environment variables are empty")
	}
}

func Start() {
	//Sanity Check
	sanityCheck()

	// Create a multiplexer
	router := mux.NewRouter()

	// wiring
	// handler -> Service -> repository (dependency injection)
	//ch := CustomerHandlers{service.NewCustomerService(domain.NewCustomerRepositoryStub())}
	pool := getDbClient()

	// Repositories
	customerRepositoryDB := domain.NewCustomerRepositoryDb(pool)
	accountRepositoryDB := domain.NewAccountRepositoryDb(pool)

	// Services
	customerService := service.NewCustomerService(customerRepositoryDB)
	accountService := service.NewAccountService(accountRepositoryDB)

	ch := CustomerHandlers{customerService}
	ah := AccountHandlers{accountService}

	// Routes

	router.HandleFunc("/customers", ch.getAllCustomers).
		Methods(http.MethodGet).
		Name("GetAllCustomers")
	router.HandleFunc("/customers/{customer_id:[0-9]+}", ch.getCustomer).
		Methods(http.MethodGet).
		Name("GetCustomer")
	router.HandleFunc("/customers/{customer_id:[0-9]+}/account", ah.createAccount).
		Methods(http.MethodPost).
		Name("NewAccount")
	router.HandleFunc("/customers/{customer_id:[0-9]+}/account/{account_id:[0-9]+}", ah.MakeTransaction).
		Methods(http.MethodPost).
		Name("NewTransaction")

	// Middlewares

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// before

			next.ServeHTTP(w, r)
			// after

		})
	})

	// Starting Server
	// OS env: SERVER_ADDRESS=localhost SERVER_PORT=8080 go run main.go
	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")
	err := http.ListenAndServe(fmt.Sprintf("%v:%v", address, port), router)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getDbClient() *pgxpool.Pool {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbAddress := os.Getenv("DB_ADDRESS")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSource := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", dbUser, dbPassword, dbAddress, dbPort, dbName)

	//DB_USER=hello_fastapi DB_PASSWORD=hello_fastapi DB_ADDRESS=localhost DB_PORT=5432 DB_NAME=banking

	pool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Println("Error connecting to DB" + err.Error())
		panic(err)
	}
	return pool
}
