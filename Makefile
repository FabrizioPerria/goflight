setup-env:
	@go get go.mongodb.org/mongo-driver/mongo
	@go get github.com/gofiber/fiber/v2
	@go install github.com/cosmtrek/air@latest
	@open -a docker && while ! docker info > /dev/null 2>&1; do sleep 1 ; done
	@docker image pull mongo:latest
	@docker network create mongonet

key:
	rm -rf ./rs_keyfile
	openssl rand -base64 756 > ./rs_keyfile
	chmod 0400 ./rs_keyfile

mongo: clean-container key
	@docker compose up -d --wait 

new-db:
	@go run scripts/db-starter.go

build:
	@go build -o bin/api .

run: build clean mongo
	@./bin/api

test:
	@go test -v ./... --count=1

clean-server:
	@lsof -i :5001 -t | grep '[0-9]*' | xargs kill -9 || true

clean-container:
	@docker-compose down

clean: clean-server clean-container

air: mongo
	@air
