package command

import "github.com/google/uuid"

type CreatePost struct {
	Id          string    `json:"id" binding:"required"`
	Text        string    `json:"text"`
	ImageBase64 string    `json:"image_base64"`
	UserId      uuid.UUID `json:"-"`
}
