package server

import (
	"Forum/internal/store"
	"Forum/internal/store/sqlite"
	"log"
	"net/http"
)

// main server structure.
type server struct {
	store  store.Store
	router *http.ServeMux
	logger *log.Logger
}

// new server instance with the given store.
func NewServer(store store.Store) *server {
	return &server{
		store:  store,
		router: &http.ServeMux{},
		logger: log.Default(),
	}
}

// ServeHTTP implements the http.Handler interface for the server.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Start initializes and starts the server with the provided configuration.
func Start(con Config) error {
	// Initialize the database
	db, err := InitDB(con)
	if err != nil {
		return err
	}

	// Create a store using SQLite implementation
	store := sqlite.NewSQL(db)

	// Create a new server instance
	server := NewServer(store)

	// Register endpoint handlers
	server.HandlePaths()

	// Start the server and log the address
	log.Println("Starting server: http://localhost:8080")
	return http.ListenAndServe(con.Port, server)
}
