# Varibles
GOCMD := go
AIRCMD := air

NAME  ?=
FILE ?= 

dev:
	$(AIRCMD) 

run:
	$(GOCMD) run main.go

migrate-create:
	$(GOCMD) run cmd/migrations/main.go create $(NAME)

migrate-up:
	$(GOCMD) run cmd/migrations/main.go up