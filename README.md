### Karma8 competition prject

#### TODO

- [ ] allocations release after failures
- [ ] cleanup after tests

#### API description

Store file endpoint:

```http://localhost:8080/api/v1/store```

Parameters: "key" - name of object to store

Download file endpoint:

```http://localhost:8080/api/v1/download```

Parameters: "key" - name of object to download


#### Running

```go run cmd/main.go```

#### Testing

```go test ./...```