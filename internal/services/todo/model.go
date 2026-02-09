package todo

import (
	pb "github.com/mmycin/goforge/proto/todo/gen"
	"strconv"
	"time"
)

// Todo represents a single todo item in the database
type Todo struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"size:255;not null" json:"description"`
	Completed   bool      `gorm:"default:false" json:"completed"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (t *Todo) ToProto() *pb.Todo {
	return &pb.Todo{
		Id:          strconv.FormatUint(uint64(t.ID), 10),
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
	}
}

func (t *Todo) ToModel() *Todo {
	return &Todo{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
	}
}
