.PHONY: make-env generate server-up server-down server-reset run-app test

make-env:
	@if [ ! -f .env ]; then \
		cp .env.sample .env; \
	fi  

generate:
	go tool templ generate
	sqlc generate

#si falla al iniciar, tambien limpia los contenedores
server-up: make-env generate
	@docker compose up -d --build --wait db api || { \
		echo "Falló el arranque. Últimos logs:"; \
		docker compose logs --tail=50 db api; \
		docker compose down -v; \
		exit 1; \
	}

server-down:
	docker compose down

server-reset:
	docker compose down -v

run-app:
	@trap '$(MAKE) --no-print-directory server-reset' EXIT; \
	trap 'exit 0' INT TERM; \
	$(MAKE) --no-print-directory server-reset && \
	$(MAKE) --no-print-directory server-up && \
	echo "Aplicación disponible en http://localhost:8080 (Ctrl+C para detener y limpiar)" && \
	docker compose logs -f api

test:
	@trap '$(MAKE) --no-print-directory server-reset' EXIT; \
	trap 'exit 130' INT TERM; \
	$(MAKE) --no-print-directory server-reset && \
	$(MAKE) --no-print-directory server-up && \
	docker compose --profile test run --rm hurl