package main

import (
	"os"

	"github.com/Blackmamoth/vault/cmd/api"
	"github.com/Blackmamoth/vault/config"
)

func main() {
	apiServer := api.NewAPIServer(config.AppConfig.APP_PORT, nil)
	if err := apiServer.Run(); err != nil {
		config.Logger.ERROR(err.Error())
		os.Exit(1)
	}
}
