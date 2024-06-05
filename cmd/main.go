package main

import (
	"os"

	"github.com/Blackmamoth/fileforte/cmd/api"
	"github.com/Blackmamoth/fileforte/config"
	"github.com/Blackmamoth/fileforte/db"
)

func main() {

	apiServer := api.NewAPIServer(config.AppConfig.APP_PORT, db.DB)

	if err := apiServer.Run(); err != nil {
		config.Logger.ERROR(err.Error())
		os.Exit(1)
	}
}
