package usecase

import (
	"github.com/ciazhar/zharapi/gen/template/data"
	"github.com/ciazhar/zharapi/gen/template/initialize"
	"github.com/ciazhar/zharapi/gen/template/initialize/db"
	config "github.com/ciazhar/zharapi/gen/template/initialize/env"
	error2 "github.com/ciazhar/zharapi/gen/template/initialize/error"
	"github.com/ciazhar/zharapi/gen/template/initialize/go_module"
	"github.com/ciazhar/zharapi/gen/template/initialize/logger"
	"github.com/ciazhar/zharapi/gen/template/initialize/middleware"
	"github.com/ciazhar/zharapi/gen/template/initialize/rest"
	string2 "github.com/ciazhar/zharapi/gen/template/initialize/string"
	"github.com/ciazhar/zharapi/gen/template/initialize/validator"
)

func Init(d data.Data, funcMap map[string]interface{}) {
	go_module.InitGoModule(d, funcMap)
	config.InitConfigFile(d, funcMap)
	initialize.InitApp(d, funcMap)
	config.InitEnv(d, funcMap)
	initialize.InitMain(d, funcMap)
	initialize.InitGRPCMain(d, funcMap)
	middleware.InitAuth(d, funcMap)
	middleware.InitCORS(d, funcMap)
	initialize.InitGateway(d, funcMap)
	error2.InitError(d, funcMap)
	logger.InitLogger(d, funcMap)
	db.InitPG(d, funcMap)
	validator.InitValidator(d, funcMap)
	rest.InitGinRequest(d, funcMap)
	rest.InitRequest(d, funcMap)
	rest.InitResponse(d, funcMap)
	string2.InitString(d, funcMap)
	go_module.GoModTidy()
}
