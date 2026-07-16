package todo

// Topic constants for todo events.
// Use these when registering listeners and dispatching events.
const (
	EventTodoCreated = "todo.created"
	EventTodoUpdated = "todo.updated"
	EventTodoDeleted = "todo.deleted"
)

// TodoCreated is fired after a todo is successfully inserted.
type TodoCreated struct {
	ID    uint
	Title string
}

// TodoUpdated is fired after a todo's fields are changed.
type TodoUpdated struct {
	ID    uint
	Title string
}

// TodoDeleted is fired after a todo is removed.
type TodoDeleted struct {
	ID uint
}
