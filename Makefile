.PHONY: run/api 
run/api:
	go run ./cmd/api

.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags='-s -X main.buildTime=${current_time}' -o=./bin/api ./cmd/api


.PHONY: run/prod
run/prod:
	@echo 'Starting'
	./api -db-dsn=${DSN} -port=${API_PORT}
