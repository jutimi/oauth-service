# Varibles
GOCMD := go
AIRCMD := air

FILE  ?= ""

dev:
	$(AIRCMD) 

run:
	$(GOCMD) run main.go

migrate-create:
	$(GOCMD) run cmd/migrations/main.go create $(FILE)

migrate-up:
	$(GOCMD) run cmd/migrations/main.go up $(FILE)

migrate-down:
	$(GOCMD) run cmd/migrations/main.go down $(FILE)