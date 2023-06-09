package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/zaurbekzet/todo-server/internal/model"
	"github.com/zaurbekzet/todo-server/internal/store"
)

// IDCtxKey — ключ для передачи ID задачи в контексте.
type IDCtxKey struct{}

type Server struct {
	store *store.Store
}

func New(store *store.Store) *Server {
	return &Server{
		store: store,
	}
}

func (s *Server) CreateTodo(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var todo model.Todo
	if err := json.Unmarshal(b, &todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.store.CreateTodo(todo.Name, todo.Priority)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetTodo(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(IDCtxKey{}).(uint32)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	todo, err := s.store.GetTodo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	b, err := json.Marshal(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func (s *Server) GetAll(w http.ResponseWriter, r *http.Request) {
	const (
		doneKey     = "done"
		priorityKey = "priority"
	)

	var (
		filters []store.Filter
		query   = r.URL.Query()
	)

	if query.Has(doneKey) {
		done, err := strconv.ParseBool(query.Get(doneKey))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filters = append(filters, store.WithDone(done))
	}

	if query.Has(priorityKey) {
		priority, err := strconv.ParseUint(query.Get(priorityKey), 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		filters = append(filters, store.WithPriority(uint8(priority)))
	}

	todos := s.store.GetAll(filters...)

	b, err := json.Marshal(todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func (s *Server) ToggleTodo(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(IDCtxKey{}).(uint32)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	todo, err := s.store.ToggleTodo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	b, err := json.Marshal(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func (s *Server) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(IDCtxKey{}).(uint32)
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	err := s.store.DeleteTodo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
