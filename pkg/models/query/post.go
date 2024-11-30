package query

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id          uuid.UUID `json:"id" binding:"required"`
	Text        string    `json:"text"`
	ImageBase64 string    `json:"image_base64"`
	ChangeDate  time.Time `json:"change_date"`
	UserId      uuid.UUID `json:"user_id"`
}
