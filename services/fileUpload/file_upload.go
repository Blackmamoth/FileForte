package fileupload

import (
	"fmt"
	"net/http"

	authMiddleware "github.com/Blackmamoth/fileforte/middleware/auth"
	fileuploadMiddleware "github.com/Blackmamoth/fileforte/middleware/file_upload"
	fileModel "github.com/Blackmamoth/fileforte/models/file"
	"github.com/Blackmamoth/fileforte/types"
	"github.com/Blackmamoth/fileforte/utils"
	"github.com/go-chi/chi/v5"
)

func FileUploadHandler() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/upload", authMiddleware.WithJWTToken(fileuploadMiddleware.CheckFilePayload(fileuploadMiddleware.RenameFiles(uploadFile))))

	return router
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	var fileHeaders []*types.FileHeader = r.Context().Value(fileuploadMiddleware.FileKey).([]*types.FileHeader)

	if fileHeaders == nil {
		utils.SendAPIErrorResponse(w, http.StatusBadRequest, fmt.Errorf("cannot access file headers"))
	}

	for _, fileHeader := range fileHeaders {
		err := utils.WriteFile(fileHeader)
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
		_, err = fileModel.AddFile(types.AppFile{
			FileName:     fileHeader.NewName,
			OriginalName: fileHeader.Filename,
			FileSize:     fileHeader.FileHeader.Size,
			ContentType:  fileHeader.FileHeader.Header["Content-Type"][0],
		})

		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	utils.SendAPIResponse(w, http.StatusOK, "File upload successful.", nil)
}
