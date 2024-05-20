.PHONY: all clean build run

BINARY_NAME=masjid_namaz_timing

all: build

build:
	go build -o $(BINARY_NAME) main.go

run:
	go run main.go serve

clean:
	go clean
	rm -f $(BINARY_NAME)

docker_compose_run:
	docker-compose up -d

docker_compose_stop:
	docker-compose down

# export_env:
# 	.\scripts\export_config.sh .\config\config.yaml