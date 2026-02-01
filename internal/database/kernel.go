package database

import (
	"github.com/mmycin/goforge/internal/services/todo"
)

func Model() []any {
	return []any{
		&todo.Todo{},
	}
}
