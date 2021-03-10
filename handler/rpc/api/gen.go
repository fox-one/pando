package api

//go:generate protoc -I . pando.proto --twirp_out=. --go_out=.
//go:generate protoc-go-inject-tag -input=./pando.pb.go
