
all:  easyjson


# run this on files with structs tagged with //easyjson:json
# and this will generate json marshal code
easyjson:
	easyjson -all internal/gitlab/mergerequest.go
	easyjson -all internal/gitlab/requestmap.go
	easyjson -all internal/gitlab/eventmodel.go
