install:
	go install \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
    github.com/golang/protobuf/protoc-gen-go
	github.com/fatih/gomodifytags

example:
	make init package="github.com/ciazhar/zharapi"
	make generate package="github.com/ciazhar/zharapi" name="List"
	zharapi -app="init" -package="github.com/ciazhar/zharapi"
	zharapi -app="module" -package="github.com/ciazhar/zharapi" -name="List"