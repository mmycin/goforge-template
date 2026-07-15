package todo

import (
	"context"
	"fmt"

	"github.com/mmycin/goforge/core/database"
)

type TodoService struct {
	db *database.Database
}

func NewTodoService(db *database.Database) *TodoService {
	return &TodoService{db: db}
}

func (s *TodoService) GetAll(ctx context.Context) ([]*Todo, error) {
	gormDB, err := s.db.GetGorm()
	if err != nil {
		return nil, fmt.Errorf("database unavailable: %w", err)
	}

	var todos []*Todo
	if err := gormDB.WithContext(ctx).Find(&todos).Error; err != nil {
		return nil, fmt.Errorf("query todos: %w", err)
	}
	return todos, nil
}

func (s *TodoService) GetByID(ctx context.Context, id uint) (*Todo, error) {
	gormDB, err := s.db.GetGorm()
	if err != nil {
		return nil, fmt.Errorf("database unavailable: %w", err)
	}

	var todo Todo
	if err := gormDB.WithContext(ctx).First(&todo, id).Error; err != nil {
		return nil, fmt.Errorf("find todo %d: %w", id, err)
	}
	return &todo, nil
}

func (s *TodoService) Create(ctx context.Context, title, description string) (*Todo, error) {
	gormDB, err := s.db.GetGorm()
	if err != nil {
		return nil, fmt.Errorf("database unavailable: %w", err)
	}

	todo := &Todo{
		Title:       title,
		Description: description,
	}
	if err := gormDB.WithContext(ctx).Create(todo).Error; err != nil {
		return nil, fmt.Errorf("create todo: %w", err)
	}
	return todo, nil
}