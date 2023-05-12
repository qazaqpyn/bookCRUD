package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefreshSession struct {
	Id        primitive.ObjectID `json:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"user_id"`
	Token     string             `json:"token"`
	ExpiresAt time.Time
}
