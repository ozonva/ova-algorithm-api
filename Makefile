build:
	go build -o bin/main cmd/ova-algorithm-api/main.go

run:
	go run cmd/ova-algorithm-api/main.go

generate:
	mockgen -source=internal/repo/repo.go > internal/mock_repo/mock_repo.go
	mockgen -source=internal/flusher/flusher.go > internal/mock_flusher/mock_flusher.go

test:
	go test -race ./...
