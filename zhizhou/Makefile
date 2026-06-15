# 知舟 Makefile
.PHONY: server-run web-run migrate db-up init

init:
	cd server && go mod tidy
	cd web && npm install

server-run:
	cd server && go run cmd/api/main.go

web-run:
	cd web && npm run dev

db-up:
	docker compose -f deploy/docker-compose.yml up -d db

migrate:
	migrate -path server/migrations -database ${DATABASE_URL} up