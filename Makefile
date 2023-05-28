

.PHONY: migration-init
migration-init:
	goose -dir "./migrations" create "" sql

.PHONY: migration-up
migration-up:
	goose -dir "./migrations" postgres "user=postgres password=postgres dbname=postgres host=localhost port=5432" up

.PHONY: migration-down
migration-down:
	goose -dir "./migrations" postgres "user=postgres password=postgres dbname=postgres host=localhost port=5432" down