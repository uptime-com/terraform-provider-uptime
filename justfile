set dotenv-load
set positional-arguments

default:
    @just --list

generate:
  go generate ./...

snapshot:
  go run github.com/goreleaser/goreleaser@v1 build --snapshot --single-target --clean

@testacc what:
   go test -v ./internal/provider -run ${1:-"^TestAcc[A-Z]"}
