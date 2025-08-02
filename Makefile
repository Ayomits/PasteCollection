migrations_create:
	@read -p "Enter migration name: " name; \
	./backend/cmd/goose/goose -s create $$name sql

migrations_migrate:
	./backend/cmd/goose/goose up

migrations_down:
	./backend/cmd/goose/goose down
