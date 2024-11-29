package main

import (
	_ "github.com/tclutin/classflow-api/docs"
	"github.com/tclutin/classflow-api/internal/app"
	"golang.org/x/net/context"
)

//	@title			ClassFlow API
//	@version		1.0
//	@description	This is t

//	@host						localhost:8080
//	@BasePath					/api/v1
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Description for what is this security definition being used

func main() {

	app.NewApp().Run(context.Background())
}
