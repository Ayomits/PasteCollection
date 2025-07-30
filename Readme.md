# Как запустить ?

1. Установи гошку
   https://go.dev/dl/

2. Перейди в папку бекенда

```bash
cd backend
```

3. Установи зависимости

```bash
go mod vendor
```

4. Подготовь постгрес

Лёгкий способ:

```bash
docker compose up -d
```

Потруднее:

```
https://www.postgresql.org/download/
```

5. Заполни env

```.env
SECRET_API_TOKEN="super-secret-key"

GOOSE_DRIVER="postgres"
GOOSE_DBSTRING="postgres://postgres:postgres@localhost:5432"
GOOSE_MIGRATION_DIR="./migrations"
GOOSE_TABLE="migrations"
```

6. Выполни миграции

```bash
./cmd/goose/goose up
```

7. Запусти бекенд

```bash
go run cmd/app/main.go
```


8. Документация
Будет находится по маршруту `/api/docs`
