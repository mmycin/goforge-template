atlas migrate diff --env gorm init_schema
sed -i 's/`//g' internal/database/migrations/*.sql
sqlc generate