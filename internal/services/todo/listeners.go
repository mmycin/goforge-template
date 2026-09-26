package todo

import (
	"context"
	"fmt"

	"github.com/mmycin/goforge/core/database"
	coreevents "github.com/mmycin/goforge/core/events"
)

// AuditOnCreated writes an audit row whenever a todo is created.
type AuditOnCreated struct {
	DB *database.Database
}

func (l *AuditOnCreated) Handle(ctx context.Context, e TodoCreated) error {
	gormDB, err := l.DB.GetGorm()
	if err != nil {
		return fmt.Errorf("audit: database unavailable: %w", err)
	}
	return gormDB.WithContext(ctx).Create(&AuditLog{
		Action:     "created",
		ResourceID: e.ID,
		Note:       e.Title,
	}).Error
}

// AuditOnUpdated writes an audit row whenever a todo is updated.
type AuditOnUpdated struct {
	DB *database.Database
}

func (l *AuditOnUpdated) Handle(ctx context.Context, e TodoUpdated) error {
	gormDB, err := l.DB.GetGorm()
	if err != nil {
		return fmt.Errorf("audit: database unavailable: %w", err)
	}
	return gormDB.WithContext(ctx).Create(&AuditLog{
		Action:     "updated",
		ResourceID: e.ID,
		Note:       e.Title,
	}).Error
}

// AuditOnDeleted writes an audit row whenever a todo is deleted.
type AuditOnDeleted struct {
	DB *database.Database
}

func (l *AuditOnDeleted) Handle(ctx context.Context, e TodoDeleted) error {
	gormDB, err := l.DB.GetGorm()
	if err != nil {
		return fmt.Errorf("audit: database unavailable: %w", err)
	}
	return gormDB.WithContext(ctx).Create(&AuditLog{
		Action:     "deleted",
		ResourceID: e.ID,
	}).Error
}

// Compile-time assertions: ensure all listeners satisfy the Listener[T] interface.
var (
	_ coreevents.Listener[TodoCreated] = (*AuditOnCreated)(nil)
	_ coreevents.Listener[TodoUpdated] = (*AuditOnUpdated)(nil)
	_ coreevents.Listener[TodoDeleted] = (*AuditOnDeleted)(nil)
)
