-- Create "todos" table
CREATE TABLE todos (
  id integer NULL PRIMARY KEY AUTOINCREMENT,
  title text NOT NULL,
  description text NOT NULL,
  completed numeric NULL DEFAULT false,
  created_at datetime NULL,
  updated_at datetime NULL
);
