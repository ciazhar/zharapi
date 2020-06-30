# ZHARAPI

> Golang HTTP and GRPC Project Generator. This project using `Makefile` executing golang code to generate project.  

## Requirement
- Golang
- Makefile

## Installation
```bash
$ make install
```

## Getting Started
```bash
$ make init package="github.com/ciazhar/zharapi"
```

## Create a new module 
```bash
$ make generate package="github.com/ciazhar/zharapi" name="List"
```
And you should register generated services to the `cmd/main.go` instance:
```diff
func InitHTTP(application *app.Application) error {
	//config router api
	router := gin.New()
+   list.InitHTTP(router,"/list", application)
-   //TODO

	//middleware
	router.Use(gin.Recovery())
	router.Use(cors.Default())
	router.Use(logger.SetLogger())

	//run
	log.Info().Caller().Msg("Running in port : " + application.Env.Get("port"))
	return router.Run(":" + application.Env.Get("port"))
}
```

or `cmd/grpc-main.go` instance:
```diff
func InitGRPC(application *app.Application) error {

	address := application.Env.Get("grpc.address")
	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	//init grpc server
	s := grpc.NewServer()

	//init client
+   list.InitGRPC(s,application)
-   //TODO

	//serve grpc server
	log.Info().Caller().Msg("Running GRPC in port : " + address)
	return s.Serve(lis)
}
```

## Start Server
- Running HTTP Server 
```bash
$ make run
```
- Running GRPC Server
```bash
$ make run-grpc
```