build:
	go fmt *.go
	go build
	go install
test:
	mkdir -p ~/Programming/Go
	go get github.com/thoj/go-ircevent
	go get golang.org/x/net/html
	go fmt *.go
	go vet *.go
	go test -v
