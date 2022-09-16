.PHONY: test fmt vet run image addmigration migrateup

ALL_PACKAGES=$(shell go list ./...)

addmigration:
	migrate create -ext sql -dir resources/migrations $(file)

test:
	go test $(ALL_PACKAGES)

fmt:
	go fmt $(ALL_PACKAGES)

vet:
	go vet $(ALL_PACKAGES)

migrateup:
	go run main.go migrate

run: migrateup
	go run main.go startserver

image:
	docker build -t url-shortener:latest .
