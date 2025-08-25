init:
	make down
	make init_var
	make build_docker
	make up
	make update_app
	make migrations_migrate
	make restart

update_app:
	docker exec -it pastcollection-bot pnpm install
	docker exec -it pastcollection-backend go mod vendor
	docker exec -it pastcollection-backend go install github.com/pressly/goose/v3/cmd/goose@latest

init_var:
	mkdir -p ./var/storage/postgres_data

down:
	docker compose down --remove-orphans

up:
	docker compose up -d

restart:
	docker compose restart

build_docker: init_var
	docker compose build

migrations_create:
	@read -p "Enter migration name: " name; \
	docker exec -it pastcollection-backend goose -s create $$name sql

migrations_migrate:
	docker exec -it pastcollection-backend goose up

migrations_down:
	docker exec -it pastcollection-backend goose down

console_backend:
	docker exec -it pastcollection-backend sh

console_bot:
	docker exec -it pastcollection-bot sh

logs_bot:
	docker compose logs pastcollection-bot -f

logs_backend:
	docker compose logs pastcollection-backend -f
