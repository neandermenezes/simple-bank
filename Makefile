container:
	docker run --name simple-bank -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres

createdb:
	docker exec -it simple-bank createdb --username=root postgres

dropdb: 
	docker exec -it simple-bank dropdb --username=root postgres

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/postgres?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/postgres?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -destination ./db/mock/store.go -package mock_db simplebank/db/sqlc Store

.PHONY: createdb dropdb migrateup migratedown sqlc test server mock