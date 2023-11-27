package utils

import (
	"time"

	"example.com/dynamicWordpressBuilding/internal/model"
	"github.com/google/uuid"
)

type Maker interface {
	CreateToken(id uuid.UUID, username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*model.Payload, error)
}
