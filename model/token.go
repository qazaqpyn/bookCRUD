package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshSession struct {
	Id        primitive.ObjectID `json:"_id,omitempty"`
	SessionId primitive.ObjectID `json:"session_id"`
	ExpiresAt time.Time
}
