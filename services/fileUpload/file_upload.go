package fileupload

import (
	"net/http"

	authMiddleware "github.com/Blackmamoth/fileforte/middleware/auth"
	fileuploadMiddleware "github.com/Blackmamoth/fileforte/middleware/file_upload"
	"github.com/Blackmamoth/fileforte/utils"
	"github.com/go-chi/chi/v5"
)

func FileUploadHandler() *chi.Mux {
	router := chi.NewRouter()

	router.Post("/upload", authMiddleware.WithJWTToken(fileuploadMiddleware.CheckFilePayload(fileuploadMiddleware.RenameFiles(uploadFile))))

	return router
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"ctx":     r.Context().Value(fileuploadMiddleware.FileKey),
		"message": "File upload successful.",
	}
	utils.SendAPIResponse(w, http.StatusOK, data, nil)
}
