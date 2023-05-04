package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Book struct {
	Id          primitive.ObjectID `json:"_id,omitempty"`
	Name        string             `json:"name" binding:"required"`
	Author      string             `json:"author"`
	Description string             `json:"description"`
	Year        time.Time          `json:"year"`
}
