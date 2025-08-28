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
	docker exec -it pastcollection-backend go install github.com/pressly/goose/v3/cmd/goose@latest
	docker exec -it pastcollection-backend go install github.com/swaggo/swag/cmd/swag@latest
	docker exec -it pastcollection-backend go mod vendor

swag_init:
	docker exec -it pastcollection-backend ./scripts/swag.sh

swag_fmt:
	docker exec -it pastcollection-backend ./scripts/swag_fmt.sh

init_var:
	mkdir -p ./var/storage/postgres_data

down:
	docker compose down --remove-orphans

up:
	docker compose up -d

restart:
	docker compose restart

restart_backend:
	docker compose restart pastcollection-backend

restart_bot:
	docker compose restart pastcollection-bot

restart_postgres:
	docker compose restart pastcollection-postgres

build_docker: init_var
	docker compose build --no-cache

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

console_postgres:
	docker exec -it pastcollection-postgres sh

cron_restore:
	TIMESTAMP=$$(date +%s%3N); \
	filename="$${TIMESTAMP}_backup"; \
	$(MAKE) name="$$filename" create_dump_auto; \
	mv "$$filename.gz" ~/PasteCollectionDumps/; \
	cd ~/PasteCollectionDumps && \
	git add "$$filename.gz" && \
	git commit -m "chore: dump $$TIMESTAMP" && \
	git push origin main

create_dump:
	@read -p "Dump name: " name; \
	make name=$$name create_dump_auto

restore_dump:
	@read -p "Dump name: " path; \
	make path=$$path restore_dump_auto

create_dump_auto:
	docker exec pastcollection-postgres pg_dump -U postgres -h localhost -p 5432 -F c -d postgres | gzip > $$name.gz; \
	echo "Dump created: $$name.gz"

restore_dump_auto:
	if [ ! -f "$$path" ]; then \
		echo "Error: File $$path not found"; \
		exit 1; \
	fi; \
	docker cp "$$path" pastcollection-postgres:/tmp/backup.gz; \
	docker exec pastcollection-postgres sh -c "gunzip -c /tmp/backup.gz | pg_restore -U postgres -h localhost -p 5432 -d postgres --clean --if-exists"; \
	docker exec pastcollection-postgres rm -f /tmp/backup.gz; \
	echo "Dump restored successfully from $$path"

logs_bot:
	docker compose logs pastcollection-bot -f

logs_backend:
	docker compose logs pastcollection-backend -f
