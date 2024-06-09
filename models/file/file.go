package fileModel

import (
	"database/sql"

	"github.com/Blackmamoth/fileforte/db"
	"github.com/Blackmamoth/fileforte/types"
)

func AddFile(appFile types.AppFile) (sql.Result, error) {
	result, err := db.DB.Exec("INSERT INTO files(`fileName`, `originalName`, `fileSize`, `contentType`) VALUES(?, ?, ?, ?)", appFile.FileName, appFile.OriginalName, appFile.FileSize, appFile.ContentType)
	return result, err
}
