.PHONY: make-env generate server-up server-down server-reset

make-env:
	@if [ ! -f .env ]; then \
		cp .env.sample .env; \
	fi  

generate:
	go tool templ generate
	sqlc generate

server-up: make-env generate
	@docker compose up -d --build --wait db api || { \
		echo "Falló el arranque. Últimos logs:"; \
		docker compose logs --tail=50 db api; \
		docker compose down; \
		exit 1; \
	}

server-down:
	docker compose down

server-reset:
	docker compose down -v
