package config

import (
	"fmt"

	golog "github.com/blackmamoth/GoLog"
)

func initializeLogger() golog.Logger {
	Logger := golog.New()
	if AppConfig.ENVIRONMENT == "DEVELOPMENT" {
		Logger.Set_Log_Level(golog.LOG_LEVEL_DEBUG)
	}
	Logger.Set_Log_Stream(golog.LOG_STREAM_MULTIPLE)
	Logger.Set_File_Name(fmt.Sprintf("%s%s", AppConfig.LOG_FILE_PATH, AppConfig.LOG_FILE_NAME))
	Logger.With_Emoji(true)
	Logger.Set_Log_Format("[%(asctime)] %(levelname) - %(message)")
	return Logger
}

var Logger = initializeLogger()
