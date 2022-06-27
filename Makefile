.PHONY: addmigration

addmigration:
	migrate create -ext sql -dir resources/migrations $(file)
