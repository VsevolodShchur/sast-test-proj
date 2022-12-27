all: build down up

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

migrate-up:
	docker-compose exec db psql -U postgres -f /migrations/up.sql

migrate-down:
	docker-compose exec db psql -U postgres -f /migrations/down.sql

.PHONY: build up down
