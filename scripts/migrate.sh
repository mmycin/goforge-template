atlas migrate diff --env gorm user_schema_2
sed -i 's/`//g' internal/database/migrations/*.sql
sqlc generate