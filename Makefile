.PHONY: up down restart logs clean fmt

up:
	docker compose up -d

down:
	docker compose down

restart:
	docker compose down
	docker compose up -d

logs:
	docker compose logs -f

clean:
	docker compose down -v
	docker system prune -f

fmt:
	@echo "Formatting services..."