migrate-and-gen:
	atlas migrate diff --env gorm $(name)
	sed -i 's/`//g' internal/database/migrations/*.sql
	sqlc generate