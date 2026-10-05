.PHONY: make-env generate test

make-env:
	@if [ ! -f .env ]; then \
		cp .env.sample .env; \
	fi  

generate:
	sqlc generate

hurl-test: make-env generate
	docker compose up -d --build db api
	docker compose run --rm hurl ; \
	status=$$?; \
    docker compose down -v; \
    exit $$status