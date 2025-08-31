
all: build


GO_FILES = $(shell find . -name '*.go')

test: $(GO_FILES)
	go test -tags test ./...
	


build: build/cmr

build/cmr: $(GO_FILES) | build/
	go build -o build/cmr

build/:
	mkdir -p build

# note assumes easyjson is installed
# run this on files with structs tagged with //easyjson:json
# and this will generate json marshal code
.PHONY: build_easy_json 
build_easy_json: internal/gitlab/requestmap_easyjson.go internal/gitlab/eventmodel_easyjson.go internal/gitlab/mergerequest_easyjson.go

internal/gitlab/requestmap_easyjson.go: internal/gitlab/requestmap.go
	easyjson -all internal/gitlab/requestmap.go

internal/gitlab/eventmodel_easyjson.go: internal/gitlab/eventmodel.go
	easyjson -all internal/gitlab/eventmodel.go

internal/gitlab/mergerequest_easyjson.go: internal/gitlab/mergerequest.go
	easyjson -all internal/gitlab/mergerequest.go

# TODO tools figure out install fresh, clean all and upgrade to latest
.PHONY: install_easy_json
install_easy_json:
	@which easyjson > /dev/null 2>&1 || go install github.com/mailru/easyjson/easyjson@latest
