default:
    exit 0

generate:
  go generate ./...

snapshot:
  go run github.com/goreleaser/goreleaser@v1 build --snapshot --single-target --clean
