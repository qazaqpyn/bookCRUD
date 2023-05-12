package repository

import (
	"context"

	"github.com/qazaqpyn/bookCRUD/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokensRepository struct {
	db *mongo.Collection
}

func NewTokensRepository(db *mongo.Database, collection string) *TokensRepository {
	return &TokensRepository{
		db: db.Collection(collection),
	}
}

func (t *TokensRepository) Create(ctx context.Context, token model.RefreshSession) error {
	_, err := t.db.InsertOne(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (t *TokensRepository) Get(ctx context.Context, token string) (model.RefreshSession, error) {
	tok := new(model.RefreshSession)

	err := t.db.FindOne(ctx, bson.M{
		"token": token,
	}).Decode(tok)
	if err != nil {
		return *tok, err
	}

	return *tok, nil
}
