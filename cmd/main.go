package main

import (
	_ "github.com/tclutin/classflow-api/docs"
	"github.com/tclutin/classflow-api/internal/app"
	"golang.org/x/net/context"
)

//	@title			ClassFlow API
//	@version		1.0
//	@description	AntonioKrasava

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Use "Bearer <token>" to authenticate

func main() {
	app.NewApp().Run(context.Background())
}
