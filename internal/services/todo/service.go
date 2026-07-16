package todo

import (
	"context"
	"fmt"

	"github.com/mmycin/goforge/core/database"
	"github.com/mmycin/goforge/core/events"
)

type TodoService struct {
	db  *database.Database
	bus *events.EventBus
}

func NewTodoService(db *database.Database, bus *events.EventBus) *TodoService {
	return &TodoService{db: db, bus: bus}
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

	events.Dispatch(ctx, s.bus, EventTodoCreated, TodoCreated{
		ID:    todo.ID,
		Title: todo.Title,
	})

	return todo, nil
}

func (s *TodoService) Update(ctx context.Context, id uint, title, description string) (*Todo, error) {
	gormDB, err := s.db.GetGorm()
	if err != nil {
		return nil, fmt.Errorf("database unavailable: %w", err)
	}

	var todo Todo
	if err := gormDB.WithContext(ctx).First(&todo, id).Error; err != nil {
		return nil, fmt.Errorf("find todo %d: %w", id, err)
	}

	todo.Title = title
	todo.Description = description

	if err := gormDB.WithContext(ctx).Save(&todo).Error; err != nil {
		return nil, fmt.Errorf("update todo %d: %w", id, err)
	}

	events.Dispatch(ctx, s.bus, EventTodoUpdated, TodoUpdated{
		ID:    todo.ID,
		Title: todo.Title,
	})

	return &todo, nil
}

func (s *TodoService) Delete(ctx context.Context, id uint) error {
	gormDB, err := s.db.GetGorm()
	if err != nil {
		return fmt.Errorf("database unavailable: %w", err)
	}

	if err := gormDB.WithContext(ctx).Delete(&Todo{}, id).Error; err != nil {
		return fmt.Errorf("delete todo %d: %w", id, err)
	}

	events.Dispatch(ctx, s.bus, EventTodoDeleted, TodoDeleted{ID: id})

	return nil
}
