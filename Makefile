BINARY_NAME=mate

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux main.go
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows main.go
	go build -o ${BINARY_NAME} main.go

run:
	./${BINARY_NAME}

build_and_run: build run

clean-builds:
	[ -f ${BINARY_NAME}-darwin ] && rm ${BINARY_NAME}-darwin || true
	[ -f ${BINARY_NAME}-linux ] && rm ${BINARY_NAME}-linux || true
	[ -f ${BINARY_NAME}-windows ] && rm ${BINARY_NAME}-windows || true
	[ -f ${BINARY_NAME} ] && rm ${BINARY_NAME} || true

clean: clean-builds
	go clean
	
test:
	go test ./...

test_coverage:
	go test ./... -coverprofile=coverage.out

dep:
	go mod download

vet:
	go vet
