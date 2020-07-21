package usecase

import (
	"github.com/ciazhar/zharapi/gen/template/data"
	"github.com/ciazhar/zharapi/gen/template/initialize/go_module"
	"github.com/ciazhar/zharapi/gen/template/module"
	"github.com/ciazhar/zharapi/gen/template/module/controller"
	"github.com/ciazhar/zharapi/gen/template/module/model"
	"github.com/ciazhar/zharapi/gen/template/module/repository"
	"github.com/ciazhar/zharapi/gen/template/module/usecase"
	"github.com/ciazhar/zharapi/gen/template/module/validator"
)

func Module(d data.Data, funcMap map[string]interface{}) {
	module.InitModuleInitializer(d, funcMap)
	model.InitModel(d, funcMap)
	model.InitProto(d, funcMap)
	repository.InitPostgresRepository(d, funcMap)
	validator.InitPostgresValidator(d, funcMap)
	repository.InitMongoRepository(d, funcMap)
	usecase.InitUseCase(d, funcMap)
	controller.InitRestController(d, funcMap)
	controller.InitGRPCController(d, funcMap)
	go_module.GoModTidy()
}
