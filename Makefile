.PHONY: create-migration build test swagger-genrate start migrate-up migrate-down start-dev
create-migration:
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

swagger-genrate:
ifneq (, $(@shell ./swagger version))
	curl -L https://github.com/go-swagger/go-swagger/releases/download/v0.26.1/swagger_linux_amd64 --output swagger
	chmod +x swagger
endif
	./swagger generate spec -o ./assets/swagger.json

start:
	go run app.go api

migration-up:
	go run app.go migrate up

migration-down:
	go run app.go migrate down
start-dev:
	@nodemon --exec go run app.go api --signal SIGTERM