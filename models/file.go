package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID       string             `json:"user_id" bson:"user_id"`
	OriginalName string             `json:"original_name" bson:"original_name"`
	StoredName   string             `json:"stored_name" bson:"stored_name"`
	MimeType     string             `json:"mime_type" bson:"mime_type"`
	Size         int64              `json:"size" bson:"size"`
	Path         string             `json:"path" bson:"path"`
	UploadedAt   time.Time          `json:"uploaded_at" bson:"uploaded_at"`
}

type FileResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	OriginalName string    `json:"original_name"`
	MimeType     string    `json:"mime_type"`
	Size         int64     `json:"size"`
	URL          string    `json:"url"`
	UploadedAt   time.Time `json:"uploaded_at"`
}

func (f *File) ToResponse(baseURL string) FileResponse {
	return FileResponse{
		ID:           f.ID.Hex(),
		UserID:       f.UserID,
		OriginalName: f.OriginalName,
		MimeType:     f.MimeType,
		Size:         f.Size,
		URL:          baseURL + "/uploads/" + f.StoredName,
		UploadedAt:   f.UploadedAt,
	}
}
