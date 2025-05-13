package csv2json

//go:generate protoc -I proto/ proto/admin.proto --go_out=admin --go-grpc_out=admin --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative
