default: .go/install

pre-push: fmt lint test

fmt:
	cd ./backend && \
	golangci-lint fmt ./...

lint:
	cd ./backend && \
	golangci-lint run ./...

test:
	cd ./backend && \
	gotestsum --junitfile report.xml --format testname -- -cover -coverprofile=coverage.out -short ./...

.go/install:
	cd ./backend && \
	go install ./cmd/...

go/mod/tidy:
	cd ./backend && \
	go mod tidy

gorm/gen:
	cd ./backend && \
	go run ./cmd/gorm/main.go && \
	git add ./infrastructure/datasource/internal/db/mapper && \
	git add ./infrastructure/datasource/internal/db/entity


db/generate/%:
	migrate create -ext sql -dir /migrations -seq $(@F)

db/up:
	migrate \
      -path "database/migrations" \
      -database "postgres://${BUSINESS_DB_USER}:${BUSINESS_DB_PASSWORD}@${BUSINESS_DB_HOST}:${BUSINESS_DB_PORT}/${BUSINESS_DB_NAME}?sslmode=${BUSINESS_DB_USE_SSL_MODE}" \
      up

db/down:
	migrate \
	  -path "database/migrations" \
	  -database "postgres://${BUSINESS_DB_USER}:${BUSINESS_DB_PASSWORD}@${BUSINESS_DB_HOST}:${BUSINESS_DB_PORT}/${BUSINESS_DB_NAME}?sslmode=${BUSINESS_DB_USE_SSL_MODE}" \
	  down
