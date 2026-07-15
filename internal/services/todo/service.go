package todo

import (
	"context"

	"github.com/mmycin/goforge/internal/database"
)

type TodoService struct {
	db *database.Database
}

func NewTodoService(db *database.Database) *TodoService {
	return &TodoService{db: db}
}

func (s *TodoService) GetAll(ctx context.Context) ([]*Todo, error) {
	// Example: use GORM to fetch todos
	gormDB, err := s.db.GetGorm()
	if err != nil {
		return nil, err
	}

	var todos []*Todo
	if err := gormDB.WithContext(ctx).Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

func (s *TodoService) GetByID(ctx context.Context, id uint) (*Todo, error) {
	gormDB, err := s.db.GetGorm()
	if err != nil {
		return nil, err
	}

	var todo Todo
	if err := gormDB.WithContext(ctx).First(&todo, id).Error; err != nil {
		return nil, err
	}
	return &todo, nil
}