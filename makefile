proto_gen:
	protoc -I pkg/grpc pkg/grpc/thumbnails.proto --go_out=./pkg/grpc/ytthumbs --go_opt=paths=source_relative --go-grpc_out=./pkg/grpc/ytthumbs --go-grpc_opt=paths=source_relative