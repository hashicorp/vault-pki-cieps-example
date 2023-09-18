
.PHONY: build
build:
	go build github.com/hashicorp/vault-pki-cieps-example/cli/cieps-server

.PHONY: fmt
fmt:
	find . -name '*.go' | grep -v pb.go | grep -v vendor | xargs go run mvdan.cc/gofumpt -w

.PHONY: certs
certs:
	openssl req -x509 -newkey rsa:2048 -keyout server.key -out server.crt -sha256 -days 3650 -nodes -subj "/CN=localhost" -addext "subjectAltName = DNS:localhost"
