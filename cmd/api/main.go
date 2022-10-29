package main

import (
	"example-restful-api-server/config"
	"example-restful-api-server/server"
	"log"

	"github.com/spf13/viper"
)

// @title           Example REST API server document
// @version         0.0.1
// @description     REST API with custom JWT-based authentication system.
// @description 	Core functionality is about creating and managing storage of photos.

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth

// @securityDefinitions.apikey  JWT
// @in                          header
// @name                        Authorization
// @description                 Description for what is this security definition being used
func main() {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	app := server.NewApp()

	if err := app.Run(viper.GetString("port")); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
