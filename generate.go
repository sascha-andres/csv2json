package csv2json

//go:generate protoc -I proto/ proto/admin.proto --go_out=pb --go-grpc_out=pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative
