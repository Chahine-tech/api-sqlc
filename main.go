package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	if users == nil {
		users = []tutorial.User{}
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

func (s *Server) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	params := tutorial.UpdateUserParams{
		ID:    int32(id),
		Name:  input.Name,
		Email: input.Email,
	}

	user, err := s.store.UpdateUser(context.Background(), params)
	if err != nil {
		http.Error(w, "Unable to update user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (s *Server) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = s.store.DeleteUser(context.Background(), int32(id))
	if err != nil {
		http.Error(w, "Unable to delete user", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Delete success"})
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
	http.HandleFunc("/updateUser/", server.updateUserHandler)
	http.HandleFunc("/deleteUser/", server.deleteUserHandler)

	log.Println("Server is starting on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
