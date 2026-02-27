run:
	cd services/web && go run ./cmd/main.go

up:
	docker compose up -d --build --force-recreate --remove-orphans

down:
	docker compose down
