package store

import (
	"errors"
	"sort"
	"sync"

	"github.com/zaurbekzet/todo-server/internal/model"
)

var ErrNotFound = errors.New("todo not found")

type Store struct {
	mu     sync.RWMutex
	todos  map[uint32]model.Todo
	nextID uint32
}

func New() *Store {
	return &Store{
		todos:  make(map[uint32]model.Todo),
		nextID: 1,
	}
}

func (s *Store) CreateTodo(name string, priority uint8) uint32 {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	s.todos[id] = model.Todo{
		ID:       id,
		Name:     name,
		Priority: priority,
	}
	s.nextID++

	return id
}

func (s *Store) GetTodo(id uint32) (model.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todo, ok := s.todos[id]
	if !ok {
		return model.Todo{}, ErrNotFound
	}

	return todo, nil
}

func (s *Store) GetAll(filters ...Filter) []model.Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]model.Todo, 0)
	for _, todo := range s.todos {
		ok := true
		for _, filter := range filters {
			if !filter(todo) {
				ok = false
				break
			}
		}
		if ok {
			todos = append(todos, todo)
		}
	}

	sort.Slice(todos, func(i, j int) bool {
		return todos[i].ID < todos[j].ID
	})

	return todos
}

func (s *Store) ToggleTodo(id uint32) (model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[id]
	if !ok {
		return model.Todo{}, ErrNotFound
	}

	todo.Done = !todo.Done
	s.todos[id] = todo

	return todo, nil
}

func (s *Store) DeleteTodo(id uint32) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.todos[id]; !ok {
		return ErrNotFound
	}
	delete(s.todos, id)

	return nil
}
