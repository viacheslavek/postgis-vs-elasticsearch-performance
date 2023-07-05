all:

.PHONY: build
build:
	docker-compose build

.PHONY: buildPostGIS
buildPostGIS:
	docker compose -f docker-compose.yml up -d dbpg

