.PHONY: generate test

generate:
	sqlc generate
	
test: generate
	docker compose up -d db
	docker compose run --rm test; \
	status=$$?; \
	docker compose down; \
	exit $$status