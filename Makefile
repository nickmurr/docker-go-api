.PHONY: migrate-up
migrate-up:
	docker run --network host migrator -path=/migrations -database "postgresql://postgres:docker@localhost:7557/postgres?sslmode=disable" up

.PHONY: migrate-down
migrate-down:
	docker run --network host migrator -path=/migrations -database "postgresql://postgres:docker@localhost:7557/postgres?sslmode=disable" down -all

.PHONY: postgres-terminal
postgres-terminal:
	psql -h localhost -p 7557 -U postgres postgres

