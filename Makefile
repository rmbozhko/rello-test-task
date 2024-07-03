migratecreate:
	docker run --rm -v $(PWD)/db/migration:/migrations migrate/migrate create -ext=sql -dir=/migrations/ -seq init_schema

migrateup:
	docker run --rm -v $(shell pwd)/db/migration:/migration --network host migrate/migrate -path=/migration/ -database "postgresql://postgres:postgres@localhost:5432/rello_test_task?sslmode=disable" -verbose up

migratedown:
	docker run --rm -v $(shell pwd)/db/migration:/migration --network host migrate/migrate -path=/migration/ -database "postgresql://postgres:postgres@localhost:5432/rello_test_task?sslmode=disable" down -all

server:
	go run cmd/main.go

.PHONY: migratecreate migrateup migratedown server