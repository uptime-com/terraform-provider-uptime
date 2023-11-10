set dotenv-load
set positional-arguments

default:
    @just --list

generate:
  go generate ./...

snapshot:
  go run github.com/goreleaser/goreleaser@v1 build --snapshot --single-target --clean

test target="^Test([^A][^c][^c]).+":
   TF_ACC="" go test -v ./internal/provider -run '{{target}}'

testacc target="^TestAcc[A-Z]":
   TF_ACC=1 go test -v ./internal/provider -run '{{target}}'
