package main

import (
	"./app"
	"./config"
)

func main() {
	appConfig := config.GetConfig()

	dcApp := &app.App{}
	dcApp.Initialize(appConfig)
	dcApp.InitializeSocket()
	dcApp.Run(":9000")
}