package store

import "github.com/zaurbekzet/todo-server/internal/model"

type Filter func(model.Todo) bool

func WithDone(done bool) Filter {
	return func(t model.Todo) bool {
		return t.Done == done
	}
}

func WithPriority(priority uint8) Filter {
	return func(t model.Todo) bool {
		return t.Priority == priority
	}
}
