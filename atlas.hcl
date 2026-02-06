data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    ".",
    "loader",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url

  // dev = "sqlite://file?mode=memory&_fk=1" # SQLITE
  dev = "mysql://root:root@localhost:3306/test" # MYSQL
  // dev = "postgres://postgres:password@localhost:5432/atlas_dev?sslmode=disable" # POSTGRES

  migration {
    dir = "file://internal/database/migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"               # keep for indentation
    }
  }
}