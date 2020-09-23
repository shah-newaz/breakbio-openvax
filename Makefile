build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags netgo -a -ldflags "-X main.VERSION=0.0.1.$$BUILD_NUMBER" -o target/breakbio-openvax && GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -tags netgo -a -ldflags "-X main.VERSION=0.0.1.$$BUILD_NUMBER" -o target/breakbio-openvax.exe

lint:
	@(which golangci-lint >/dev/null)  || (echo "Installing GolangCI-Lint"  && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.30.0)
	golangci-lint run ./... --skip-dirs mocks/
