# GIN BOILERPLATE

![Docker Image](https://github.com/jutimi/oauth-service/actions/workflows/docker-image.yml/badge.svg?branch=master)

## Tech stacks

- Framework: Gin (Golang)
- Database: Postgres
- Cache: Redis
- Message queue: Kafka
- Others: gRPC

## Prerequisites

- Install make
- Install docker/docker compose
- Install golang

## CMD

- Server

  - `make run`: Run project. Default at port 8080
  - `make dev`: Run project using air (hot reload). Default at port 8080

- Migration:
