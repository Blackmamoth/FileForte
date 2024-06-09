package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/Blackmamoth/fileforte/config"
	"github.com/Blackmamoth/fileforte/types"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func SendAPIResponse(w http.ResponseWriter, status int, data any, cookie *http.Cookie) error {
	if cookie != nil {
		http.SetCookie(w, cookie)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(generateAPIResponseBody(status, data))
}

func SendAPIErrorResponse(w http.ResponseWriter, status int, err error) {
	SendAPIResponse(w, status, err.Error(), nil)
}

func generateAPIResponseBody(status int, data any) map[string]any {
	return map[string]any{"status": status, "data": data}
}

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(v)
}

func WriteFile(fileHeader *types.FileHeader) error {
	_, err := os.Stat(config.FileUploadConfig.MEDIA_UPLOAD_PATH)
	if err != nil {
		if os.IsNotExist(err) {
			err = createMediaUploadDirectory()

			if err != nil {
				return err
			}

		} else {
			return err
		}
	}

	mediaPath := config.FileUploadConfig.MEDIA_UPLOAD_PATH
	seperator := "/"
	if runtime.GOOS == "windows" {
		seperator = "\\"
	}

	file, err := os.Create(fmt.Sprintf("%s%s%s", mediaPath, seperator, fileHeader.NewName))

	if err != nil {
		return err
	}
	defer file.Close()

	fileData, err := fileHeader.FileHeader.Open()

	if err != nil {
		return err
	}

	defer fileData.Close()

	_, err = io.Copy(file, fileData)

	return err
}

func createMediaUploadDirectory() error {
	return os.MkdirAll(config.FileUploadConfig.MEDIA_UPLOAD_PATH, os.ModePerm)
}
