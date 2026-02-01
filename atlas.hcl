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

  dev = "sqlite://file?mode=memory&_fk=1"   # in-memory for diff planning; add &_fk=1 for foreign keys if you use them

  migration {
    dir = "file://internal/database/migrations"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"               # keep for indentation
    }
  }
}