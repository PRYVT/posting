package query

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id          uuid.UUID
	Text        string
	ImageBase64 string
	ChangeDate  time.Time
}
