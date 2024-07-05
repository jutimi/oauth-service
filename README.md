# GIN BOILERPLATE

## Prerequisites

- Install make
- Install docker/docker compose

## CMD

- Server

  - `make run`: Run project. Default at port 8080
  - `make dev`: Run project using air (hot reload). Default at port 8080

- Migration:

  - `make migrate-create FILE=...`: Create migration file (Ex: make migrate-create FILE=create_table_user)
  - `make migrate-up`: Run all migrations
  - `make migrate-up FILE=...`: Run specific migration
  - `make migrate-down`: Down all migrations

## Pre Deploy
