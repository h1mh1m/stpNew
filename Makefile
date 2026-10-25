.PHONY: start-frontend start-backend migrate-up migrate-down migrate-status seed-rooms test-backend test-backend-race

start-frontend:
	cd frontend-web && npm run dev

start-backend:
	cd backend && go run main.go

migrate-up:
	cd backend && go run ./cmd/migrate -up

migrate-down:
	cd backend && go run ./cmd/migrate -down

migrate-status:
	cd backend && go run ./cmd/migrate -status

seed-rooms:
	cd backend && go run ./cmd/seed

test-backend:
	cd backend && go test -p 1 ./...

test-backend-race:
	cd backend && go test -p 1 -race ./...