package handler

import (
	"github.com/google/uuid"
)

type reply struct {
	ID uuid.UUID `json:"id"`
}
