.PHONY: create_migration build test swagger_genrate start migrate_up migrate_down
create_migration:
ifneq (, $(@shell ./migrate -version))
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.11.0/migrate.linux-amd64.tar.gz | tar xvz
	mv migrate.linux-amd64 migrate
	chmod +x migrate
endif
	./migrate create -ext sql -dir database/migrations $(file_name)

build:
	go build -o=$(app_name) .

test:
	@go test -coverprofile coverage.out ./...
	@echo "=========================================================================================="
	@echo "                                TEST COVERAGE                                             "
	@echo "=========================================================================================="
	@go tool cover -func coverage.out

swagger_genrate:
ifneq (, $(@shell ./swagger version))
	curl -L https://github.com/go-swagger/go-swagger/releases/download/v0.26.1/swagger_linux_amd64 --output swagger
	chmod +x swagger
endif
	./swagger generate spec -o ./assets/swagger.json

start:
	go run app.go api

migration_up:
	go run app.go migrate up

migration_down:
	go run app.go migrate down
