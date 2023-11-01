.PHONY: test

APP_NAME = home-monitor

build:
	go build -o bin/$(APP_NAME)

docker-build-push: build
	docker build -t andreasohlen/$(APP_NAME):latest .
	docker push andreasohlen/home-monitor:latest

test:
	go test ./config/... -v

run: build
	./bin/$(APP_NAME)