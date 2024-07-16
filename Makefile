# Varibles
GOCMD := go
AIRCMD := air

FILE  ?= ""

dev:
	$(AIRCMD) 

run:
	$(GOCMD) run main.go