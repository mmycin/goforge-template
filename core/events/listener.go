package events

import "context"

// Listener handles a specific event type T.
// Implement this interface for every listener in your application.
//
// Example:
//
//	type AuditOnCreated struct{ DB *database.Database }
//
//	func (l *AuditOnCreated) Handle(ctx context.Context, e TodoCreated) error {
//	    return l.DB.Gorm.WithContext(ctx).Create(&AuditLog{...}).Error
//	}
type Listener[T any] interface {
	Handle(ctx context.Context, event T) error
}
