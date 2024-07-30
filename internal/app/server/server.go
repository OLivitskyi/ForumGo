package server

import (
	"Forum/internal/store"
	"Forum/internal/store/sqlite"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	Store     store.Store
	router    *mux.Router
	Logger    *log.Logger
	Templates *template.Template
}

func NewServer(store store.Store) *Server {
	return &Server{
		Store:     store,
		router:    mux.NewRouter(),
		Logger:    log.Default(),
		Templates: template.Must(template.ParseGlob("../../web/templates/*.html")),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) GetUserIDFromSession(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session_uuid")
	if err != nil {
		return 0, err
	}
	sessionToken := cookie.Value
	return s.Store.Session().GetUserIDFromSession(sessionToken)
}

func Start(con Config) error {
	db, err := InitDB(con)
	if err != nil {
		return err
	}

	store := sqlite.NewSQL(db)
	server := NewServer(store)
	server.HandlePaths()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(server)

	log.Println("Starting server: http://localhost:8080")

	return http.ListenAndServe(con.Port, handler)
}
