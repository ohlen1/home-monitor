.PHONY: test

APP_NAME = home-monitor

build:
	go build -o bin/$(APP_NAME)

docker-build-push:
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t andreasohlen/$(APP_NAME):latest --push .

test:
	go test ./config/... -v
	go test ./metrics/... -v

run: build
	./bin/$(APP_NAME)