build:
	go fmt *.go
	go build
	go install
test:
	go fmt *.go
	go vet *.go
	go test -v
