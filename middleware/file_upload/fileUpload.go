package fileuploadMiddleware

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/Blackmamoth/fileforte/config"
	"github.com/Blackmamoth/fileforte/types"
	"github.com/Blackmamoth/fileforte/utils"
	"github.com/google/uuid"
)

type contextKey string

const FileKey contextKey = "files"

func CheckFilePayload(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			utils.SendAPIErrorResponse(w, http.StatusUnsupportedMediaType, fmt.Errorf("content-type should be multipart/form-data"))
			return
		}

		err := r.ParseMultipartForm(config.FileUploadConfig.MAX_UPLOAD_SIZE << 20)
		if err != nil {
			utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		if r.MultipartForm.File == nil {
			utils.SendAPIErrorResponse(w, http.StatusUnprocessableEntity, fmt.Errorf("file not present in request body"))
			return
		}

		handlerFunc(w, r)
	}
}

func RenameFiles(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files := r.MultipartForm.File["files"]
		fileHeaders := make([]*types.FileHeader, len(files))
		for i, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				utils.SendAPIErrorResponse(w, http.StatusInternalServerError, err)
				return
			}
			defer file.Close()
			newFileName := fmt.Sprintf("%s.%s", uuid.New(), filepath.Ext(fileHeader.Filename))
			fileHeaders[i] = &types.FileHeader{
				OriginalName: fileHeader.Filename,
				NewName:      newFileName,
				FileHeader:   fileHeader,
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, FileKey, fileHeaders)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}
