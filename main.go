package main

import (
	"github.com/jackson6/gcic-server/app"
	"github.com/jackson6/gcic-server/config"
)

func main() {
	appConfig := config.GetConfig()

	dcApp := &app.App{}
	dcApp.Initialize(appConfig)
	dcApp.InitializeSocket()
	dcApp.Run(":9000")
}