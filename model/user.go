package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        primitive.ObjectID `json:"id"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Email     string             `json:"email"`
	Password  string             `json:"password"`
}
