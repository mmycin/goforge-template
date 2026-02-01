-- name: CreateTodo :exec
INSERT INTO todos (title, completed) VALUES (?, ?);

-- name: GetTodo :one
SELECT * FROM todos WHERE id = ?;

-- name: ListTodos :many
SELECT * FROM todos;

-- name: UpdateTodo :exec
UPDATE todos SET title = ?, completed = ? WHERE id = ?;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = ?;