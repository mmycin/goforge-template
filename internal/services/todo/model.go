package todo

import (
	"time"
)

// Todo represents a single todo item in the database.
type Todo struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"size:255;not null" json:"description"`
	Completed   bool      `gorm:"default:false" json:"completed"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (t *Todo) ToTodoModel() *Todo {
	return &Todo{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
	}
}

// AuditLog records every create/update/delete action on a todo.
type AuditLog struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Action     string    `gorm:"size:50;not null" json:"action"`
	ResourceID uint      `gorm:"not null;index" json:"resource_id"`
	Note       string    `gorm:"size:255" json:"note"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

