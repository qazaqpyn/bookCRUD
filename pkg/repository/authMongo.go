package repository

import (
	"context"

	"github.com/qazaqpyn/bookCRUD/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository struct {
	db *mongo.Collection
}

func NewAuthRepository(db *mongo.Database, collection string) *AuthRepository {
	return &AuthRepository{
		db: db.Collection(collection),
	}
}

func (r AuthRepository) CreateUser(ctx context.Context, user *model.User) error {
	_, err := r.db.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (r AuthRepository) GetUser(ctx context.Context, email, password string) (*model.User, error) {
	user := new(model.User)

	err := r.db.FindOne(ctx, bson.M{
		"email":    email,
		"password": password,
	}).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r AuthRepository) DeleteUser(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (r AuthRepository) UpdateUser(ctx context.Context, id primitive.ObjectID, user *model.User) error {
	_, err := r.db.UpdateByID(ctx, bson.M{"id": id}, user)
	return err
}
