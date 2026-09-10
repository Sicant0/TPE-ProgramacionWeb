.PHONY: make-env generate test

make-env:
	@if [ ! -f .env ]; then \
		cp .env.sample .env; \
	fi  

generate:
	sqlc generate
	
test: make-env generate
	docker compose up -d db
	docker compose run --rm test; \
	status=$$?; \
	docker compose down; \
	exit $$status