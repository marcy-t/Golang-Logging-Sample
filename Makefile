# include .env

all:
	docker compose -f ./build/docker-compose.dev.yml up -d

up:
	docker compose -f ./build/docker-compose.dev.yml up -d

down:
	docker compose -f ./build/docker-compose.dev.yml down
