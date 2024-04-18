setup-env:
	@go get go.mongodb.org/mongo-driver/mongo
	@go get github.com/gofiber/fiber/v2
	@go install github.com/cosmtrek/air@latest
	@open -a docker && while ! docker info > /dev/null 2>&1; do sleep 1 ; done
	@docker image pull mongo:latest

mongo: clean-container
	@docker run --name mongodb -d  -p 27017:27017 mongo:latest

build:
	@go build -o bin/api .

run: build clean mongo
	@./bin/api

test:
	@go test -v ./... --count=1

clean-server:
	@lsof -i :5001 -t | grep '[0-9]*' | xargs kill -9 || true

clean-container:
	@docker inspect --format '{{json .State.Running}}' mongodb 2>/dev/null | grep true && docker stop mongodb && docker rm mongodb || true

air: mongo
	@air
