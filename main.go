package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Chahine-tech/api-sqlc/tutorial"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Server struct {
	store *tutorial.Queries
}

func NewServer(store *tutorial.Queries) *Server {
	return &Server{store: store}
}

func (s *Server) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := s.store.GetUsers(context.Background())
	if err != nil {
		http.Error(w, "Unable to fetch users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func (s *Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Créez une instance de CreateUserParams
	params := tutorial.CreateUserParams{
		Name:  input.Name,
		Email: input.Email,
	}

	user, err := s.store.CreateUser(context.Background(), params)
	if err != nil {
		http.Error(w, "Unable to create user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func main() {
	connStr := "postgresql://user:password@database:5432/database" // Utilise le nom du service Docker
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	store := tutorial.New(conn)
	server := NewServer(store)

	http.HandleFunc("/users", server.getUsersHandler)
	http.HandleFunc("/createUser", server.createUserHandler)

	log.Println("Server is starting on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
