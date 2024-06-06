all: check build

check:
	go vet ./...

build:
	go build cmd/*

clean:
	find -not -path "*.git*" -type f -executable -exec rm {} \;
