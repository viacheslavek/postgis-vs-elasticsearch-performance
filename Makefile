all:

.PHONY: build
build:
	docker-compose build

.PHONY: buildPostGIS
buildPostGIS:
	docker compose -f docker-compose.yml up -d dbpg

.PHONY: buildES
buildES:
	docker compose -f docker-compose.yml up -d elasticsearch

.PHONY: run
run:
	go run ./cmd/app

