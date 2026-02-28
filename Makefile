BUILD_FLAGS=-gcflags="all=-N -l" -trimpath -mod=readonly -modcacherw

build:
	go build $(BUILD_FLAGS) -o mkrfc .

lint:
	golangci-lint run --timeout 5m

test:
	go test -v -cover ./...

install:
	go install .

check:
	go mod tidy
	git diff --exit-code go.mod go.sum
	make lint
	make test
