include .env
dev:
	~/go/bin/air --build.cmd "go build -o ./bin/go-face-detection-be ./cmd/main.go" --build.bin ./bin/go-face-detection-be
api-docs:
	~/go/bin/swag init -g ./cmd/main.go
unit-test:
	set -a && . ./.env && go test -race -v -coverprofile=profile.out ./... $(shell echo $(TEST_FLAGS)) && go tool cover -html=profile.out ; rm -f cover.out
coverage:
	@go test -covermode=count -coverprofile=count.out fmt; rm -f count.out