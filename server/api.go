package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/cors"
)

type APIServer struct {
	addr  string
	store *Storage
}

type apiHandlerFunc func(http.ResponseWriter, *http.Request) error

func makeHandler(f apiHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			WriteError(w, http.StatusBadRequest, err)
		}
	}
}

func NewAPIServer(addr string, store *Storage) *APIServer {
	return &APIServer{
		addr:  addr,
		store: store,
	}
}

func (s *APIServer) Serve() {
	r := http.NewServeMux()

	r.HandleFunc("POST /api/todo", makeHandler(s.handleCreateTodo))
	r.HandleFunc("GET /api/todos", makeHandler(s.handleGetAllTodos))
	r.HandleFunc("GET /api/todo/{id}", makeHandler(s.handleGetTodoById))
	r.HandleFunc("PUT /api/todo/{id}", makeHandler(s.handleUpdateTodo))
	r.HandleFunc("DELETE /api/todo/{id}", makeHandler(s.handleDeleteTodo))

	corsOpts := cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}

	handler := cors.New(corsOpts).Handler(r)

	http.ListenAndServe(s.addr, handler)
}

func (s *APIServer) handleCreateTodo(w http.ResponseWriter, r *http.Request) error {
	var todo Todo

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		return err
	}
	err = s.store.CreateTodo(&todo)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusCreated, todo)
	// return WriteJSON(w, http.StatusCreated, map[string]string{"message": fmt.Sprint(todo.ID)})
}

func (s *APIServer) handleGetAllTodos(w http.ResponseWriter, r *http.Request) error {
	var todos []Todo
	err := s.store.GetAllTodos(&todos)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, todos)
}

func (s *APIServer) handleGetTodoById(w http.ResponseWriter, r *http.Request) error {
	var todo Todo
	id := r.PathValue("id")
	err := s.store.GetTodoById(id, &todo)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, todo)
}

func (s *APIServer) handleUpdateTodo(w http.ResponseWriter, r *http.Request) error {
	var todo Todo
	id := r.PathValue("id")
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		return err
	}
	err = s.store.UpdateTodo(id, &todo)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]string{"message": "todo item successfully updated", "id": fmt.Sprint(id)})
}

func (s *APIServer) handleDeleteTodo(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	err := s.store.DeleteTodo(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]string{"message": "todo item successfully deleted", "id": fmt.Sprint(id)})
}
