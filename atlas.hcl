data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "./cmd/main.go",
    "loader",
  ]
}

locals {
  db_driver = getenv("DB_CONNECTION")
  db_user   = getenv("DB_USERNAME")
  db_pass   = getenv("DB_PASSWORD")
  db_host   = getenv("DB_HOST")
  db_port   = getenv("DB_PORT")
  db_name   = getenv("DB_NAME")
  db_dev_name = getenv("DB_DEV_NAME") == "" ? getenv("DB_NAME") : getenv("DB_DEV_NAME")

  # Construct dev URL based on driver
  mysql_url    = "mysql://${local.db_user}:${local.db_pass}@${local.db_host}:${local.db_port}/${local.db_dev_name}"
  postgres_url = "postgres://${local.db_user}:${local.db_pass}@${local.db_host}:${local.db_port}/${local.db_dev_name}?sslmode=disable"
  sqlite_url   = "sqlite://file?mode=memory&_fk=1"

  dev_url = local.db_driver == "mysql" ? local.mysql_url : (
            local.db_driver == "postgres" ? local.postgres_url : (
            local.db_driver == "sqlite" ? local.sqlite_url : ""
          ))
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = local.dev_url

  migration {
    dir = "file://internal/database/migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"               # keep for indentation
    }
  }
}