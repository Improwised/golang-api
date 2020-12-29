create-migrate:
	migrate create -ext sql -dir database/migrations $(file_name)

build:
	go build -o=$(app_name) .

test:
	go test -coverprofile coverage.out ./... && go tool cover -func coverage.out
