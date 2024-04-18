setup-env:
	go get go.mongodb.org/mongo-driver/mongo
	go get github.com/gofiber/fiber/v2
	open -a docker && while ! docker info > /dev/null 2>&1; do sleep 1 ; done
	docker image pull mongo:latest

docker:
	docker run --name mongodb -d mongo:latest -p 27017:27017 

build:
	go build -o bin/api .

run: build
	./bin/api

test:
	go test -v ./...
