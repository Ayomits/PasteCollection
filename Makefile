migrations_create:
	@read -p "Enter migration name: " name; \
	./cmd/goose/goose -s create $$name sql

migrations_migrate:
	./cmd/goose/goose up

migrations_down:
	./cmd/goose/goose down
