package types

import "mime/multipart"

type FileHeader struct {
	OriginalName string
	NewName      string
	*multipart.FileHeader
}
