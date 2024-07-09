migratecreate:
	docker run --rm -v $(PWD)/db/migration:/migrations migrate/migrate create -ext=sql -dir=/migrations/ -seq init_schema

migrateup:
	docker run --rm -v $(shell pwd)/db/migration:/migration --network host migrate/migrate -path=/migration/ -database "postgresql://postgres:postgres@localhost:5432/rello_test_task?sslmode=disable" -verbose up

migratedown:
	docker run --rm -v $(shell pwd)/db/migration:/migration --network host migrate/migrate -path=/migration/ -database "postgresql://postgres:postgres@localhost:5432/rello_test_task?sslmode=disable" down -all

server:
	go run cmd/main.go

sqlboiler:
	docker run --rm -it -v $(PWD)/sqlboiler.toml:/sqlboiler.toml:ro -v $(PWD)/db/models:/models:rw --network host goodwithtech/sqlboiler:latest --wipe /sqlboiler-psql --output models/gen

swag:
	${HOME}/go/bin/swag init -d cmd,api

test:
	go test -v -cover ./...

# sqlboiler2:
# docker run --rm -v $(PWD)/sqlboiler.toml:/sqlboiler.toml -v $(PWD)/db/models:/models curvegrid/sqlboiler:latest psql

.PHONY: migratecreate migrateup migratedown server sqlboiler swag test