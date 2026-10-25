.PHONY: dev start start-frontend start-backend build

start-frontend:
	cd frontend-web && npm run dev

start-backend:
	cd backend && go run main.go

dev:
	@bash -c "trap 'kill 0' SIGINT SIGTERM EXIT; (cd backend && go run main.go) & (cd frontend-web && npm run dev)"

start: dev

build:
	cd backend && go build -o backend .
	cd frontend-web && npm run build