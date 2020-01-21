.PHONY: migrate-up
migrate-up:
	docker run --network host migrator -path=/migrations -database "postgresql://postgres:docker@localhost:7557/postgres?sslmode=disable" up

.PHONY: migrate-down
migrate-down:
	docker run --network host migrator -path=/migrations -database "postgresql://postgres:docker@localhost:7557/postgres?sslmode=disable" down -all

.PHONY: pt
pt:
	psql -h localhost -p 7557 -U postgres postgres

.PHONY: test-go
test-go:
	docker exec -it docker-examples_api_1 go test -v -coverprofile=c.out -timeout 30s ./...
