-- name: GetTodo :one
SELECT * FROM todos WHERE id = ? LIMIT 1;

-- name: ListTodos :many
SELECT * FROM todos ORDER BY created_at DESC;

-- name: CreateTodo :one
INSERT INTO todos (title, description, completed, created_at, updated_at)
VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING *;

-- name: UpdateTodo :one
UPDATE todos SET title = ?, description = ?, completed = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = ?;
