all:
	go fmt *.go
	go vet *.go
	go build
	go install
test:
	# Get imports
	mkdir -p ~/Programming/Go
	go get github.com/thoj/go-ircevent
	go get golang.org/x/net/html
	go get github.com/Luminarys/gochat
	# Test primary
	go fmt *.go
	go vet *.go
	go test -v github.com/Luminarys/gochat
	# Test modules
	go fmt modules/*.go
	go vet modules/*.go
	go test -v github.com/Luminarys/gochat/modules
	cd modules
