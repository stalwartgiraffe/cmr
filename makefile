
all: build


GO_FILES = $(shell find . -name '*.go')

test: $(GO_FILES)
	go test -tags test ./...
	
build: build/cmr

build/cmr: $(GO_FILES) | build/
	go build -o build/cmr

build/:
	mkdir -p build


# run this on files with structs tagged with //easyjson:json
# and this will generate json marshal code
easyjson:
	easyjson -all internal/gitlab/mergerequest.go
	easyjson -all internal/gitlab/requestmap.go
	easyjson -all internal/gitlab/eventmodel.go
