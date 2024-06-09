package types

import "time"

type AppFile struct {
	Id           int
	FileName     string
	OriginalName string
	FileSize     int64
	ContentType  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
