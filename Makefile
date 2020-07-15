install:
	go install \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
    github.com/golang/protobuf/protoc-gen-go
	github.com/fatih/gomodifytags

init: init-gomod \
	init-config-file \
	init-app \
	init-env \
	init-main \
	init-grpc \
	init-middleware \
	init-gateway \
	init-error \
	init-logger \
	init-pg \
	init-validator \
	init-rest \
	init-string
	go mod tidy

init-gomod:
	[ -f ./go.mod ] && echo exists || go mod init $(package) ;

init-config-file:
	go run gen/template/init/config.go \
		-package=$(package);

init-app:
	mkdir -p app
	go run gen/template/init/app.go \
		-package=$(package)

init-env:
	mkdir -p common/env
	go run gen/template/init/common-template/env/env.go \
		-package=$(package)

init-main:
	mkdir -p cmd
	go run gen/template/init/main.go \
		-package=$(package);

init-grpc:
	mkdir -p cmd
	go run gen/template/init/grpc-main.go \
		-package=$(package);

init-gateway:
	mkdir -p cmd
	go run gen/template/init/gateway.go \
		-package=$(package);

init-middleware:
	mkdir -p common/middleware
	go run gen/template/init/common-template/middleware/auth.go \
		-package=$(package)
	go run gen/template/init/common-template/middleware/cors.go \
		-package=$(package)

init-error:
	mkdir -p common/error
	go run gen/template/init/common-template/error/error.go \
		-package=$(package)
	gomodifytags -file common/error/error.go -struct Error -add-tags json -w

init-logger:
	mkdir -p common/logger
	go run gen/template/init/common-template/logger/logger.go \
		-package=$(package)

init-pg:
	mkdir -p common/db/pg
	go run gen/template/init/common-template/db/pg.go \
		-package=$(package)

init-validator:
	mkdir -p common/validator
	go run gen/template/init/common-template/validator/validator.go \
		-package=$(package)

init-string:
	mkdir -p common/string
	go run gen/template/init/common-template/string/string.go \
		-package=$(package)

init-rest:
	mkdir -p common/rest
	go run gen/template/init/common-template/rest/rest.go \
		-package=$(package)
	go run gen/template/init/common-template/rest/gin_rest.go \
		-package=$(package)
	go run gen/template/init/common-template/rest/response.go \
			-package=$(package)

generate: generate-init \
	generate-model \
	generate-proto \
	generate-postgres-repository \
	generate-postgres-validator \
	generate-usecase \
	generate-rest-controller \
	generate-grpc-controller
	go mod tidy

generate-init:
	mkdir -p src/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]')
	go run gen/template/module/init.go \
		-name=$(name) \
		-package=$(package)

generate-model:
	mkdir -p src/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]')/model
	go run gen/template/module/model/model.go \
		-name=$(name) \
		-package=$(package)

generate-proto:
	mkdir -p grpc/proto
	mkdir -p grpc/generated/golang
	mkdir -p grpc/generated/swagger
	go run gen/template/module/model/proto.go \
		-name=$(name) \
		-package=$(package)
	gomodifytags -file src/$(name)/model/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]').go -struct $(name) -add-tags json -w
	gomodifytags -file src/$(name)/model/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]').go -line 6 -add-tags pg:$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]') -w
	gomodifytags -file src/$(name)/model/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]').go -line 6 -remove-tags json -w
	protoc -I/usr/local/include -I. \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway \
		--go_out=plugins=grpc:./grpc \
		--grpc-gateway_out=logtostderr=true:./grpc \
		--swagger_out=allow_merge=true,merge_file_name=global:./grpc/generated/swagger \
		grpc/proto/**

generate-postgres-repository:
	mkdir -p src/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]')/repository/postgres
	go run gen/template/module/repository/pg.go \
    	-name=$(name) \
    	-package=$(package)

generate-postgres-validator:
	mkdir -p src/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]')/validator/postgres
	go run gen/template/module/validator/pg.go \
		-name=$(name) \
		-package=$(package)

generate-usecase:
	mkdir -p src/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]')/usecase
	go run gen/template/module/usecase/usecase.go \
		-name=$(name) \
		-package=$(package)

generate-rest-controller:
	mkdir -p src/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]')/controller/rest
	go run gen/template/module/controller/rest.go \
		-name=$(name) \
		-package=$(package)

generate-grpc-controller:
	mkdir -p src/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]')/controller/grpc
	go run gen/template/module/controller/grpc.go \
		-name=$(name) \
		-package=$(package)

run:
	go run cmd/main.go

run-grpc:
	go run cmd/grpc-main.go

run-gateway:
	go run cmd/gateway.go

example:
	make init package="github.com/ciazhar/zharapi"
	make generate package="github.com/ciazhar/zharapi" name="List"
