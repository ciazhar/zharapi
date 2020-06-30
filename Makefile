install:
	go install \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
    github.com/golang/protobuf/protoc-gen-go
	github.com/betacraft/easytags

init: init-gomod \
	init-config-file \
	init-app \
	init-env \
	init-main \
	init-grpc \
	init-error \
	init-logger \
	init-pg \
	init-validator \
	init-rest
	go mod tidy

init-gomod:
	[ -f ./go.mod ] && echo exists || go mod init $(package) ;

init-config-file:
	go run gen/template/init/config.go \
		-package=$(package) > config.json ;

init-app:
	mkdir -p app
	go run gen/template/init/app.go \
		-package=$(package) > app/app.go ;

init-env:
	mkdir -p common/env
	go run gen/template/init/common-template/env.go \
		-package=$(package) > common/env/env.go

init-main:
	mkdir -p cmd
	go run gen/template/init/main.go \
		-package=$(package) > cmd/main.go ;

init-grpc:
	mkdir -p cmd
	go run gen/template/init/grpc-main.go \
		-package=$(package) > cmd/grpc-main.go ;

init-error:
	mkdir -p common/error
	go run gen/template/init/common-template/error.go \
		-package=$(package) > common/error/error.go

init-logger:
	mkdir -p common/logger
	go run gen/template/init/common-template/logger.go \
		-package=$(package) > common/logger/log.go

init-pg:
	mkdir -p common/db/pg
	go run gen/template/init/common-template/pg.go \
		-package=$(package) > common/db/pg/pg.go

init-validator:
	mkdir -p common/validator
	go run gen/template/init/common-template/validator.go \
		-package=$(package) > common/validator/validator.go

init-rest:
	mkdir -p common/rest
	go run gen/template/init/common-template/rest.go \
		-package=$(package) > common/rest/param.go
	go run gen/template/init/common-template/gin_rest.go \
		-package=$(package) > common/rest/gin_param.go

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
		-type=$(type) \
		-package=$(package) > src/$(name)/init.go

generate-model:
	mkdir -p src/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]')/model
	go run gen/template/module/model/model.go \
		-name=$(name) \
		-type=$(type) \
		-package=$(package) > src/$(name)/model/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]').go

generate-proto:
	mkdir -p grpc/proto
	mkdir -p grpc/generated/golang
	mkdir -p grpc/generated/swagger
	go run gen/template/module/model/proto.go \
		-name=$(name) \
		-type=$(type) \
		-package=$(package) > grpc/proto/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]').proto
	easytags src/$(name)/model/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]').go
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
    	-type=$(type) \
    	-package=$(package) > src/$(name)/repository/postgres/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]')_pg_repo.go

generate-postgres-validator:
	mkdir -p src/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]')/validator/postgres
	go run gen/template/module/validator/pg.go \
		-name=$(name) \
		-type=$(type) \
		-package=$(package) > src/$(name)/validator/postgres/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]')_pg_validator.go

generate-usecase:
	mkdir -p src/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]')/usecase
	go run gen/template/module/usecase/usecase.go \
		-name=$(name) \
		-type=$(type) \
		-package=$(package) > src/$(name)/usecase/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]')_usecase.go

generate-rest-controller:
	mkdir -p src/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]')/controller/rest
	go run gen/template/module/controller/rest.go \
		-name=$(name) \
		-type=$(type) \
		-package=$(package) > src/$(name)/controller/rest/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]')_rest_controller.go

generate-grpc-controller:
	mkdir -p src/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]')/controller/grpc
	go run gen/template/module/controller/grpc.go \
		-name=$(name) \
		-type=$(type) \
		-package=$(package) > src/$(name)/controller/grpc/$(shell echo '$(name)' | tr '[:upper:]' '[:lower:]')_grpc_controller.go

example:
	make init package="github.com/ciazhar/generate"
	make generate package="github.com/ciazhar/generate" name="List" type="*list.List"
