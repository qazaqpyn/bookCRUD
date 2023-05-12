package repository

import (
	"context"

	"github.com/qazaqpyn/bookCRUD/model"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Authorization interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, email, password string) (*model.User, error)
}

type Book interface {
	CreateBook(ctx context.Context, book model.Book) error
	GetBookById(ctx context.Context, id primitive.ObjectID) (*model.Book, error)
	GetAllBook(ctx context.Context) ([]*model.Book, error)
	UpdateBook(ctx context.Context, id primitive.ObjectID, book *model.Book) error
	DeleteBook(ctx context.Context, id primitive.ObjectID) error
}

type Tokens interface {
	Create(ctx context.Context, token model.RefreshSession) error
	Get(ctx context.Context, token string) (model.RefreshSession, error)
}

type Repository struct {
	Authorization
	Book
	Tokens
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db, viper.GetString("mongo.user")),
		Book:          NewBookRepository(db, viper.GetString("mongo.book")),
		Tokens:        NewTokensRepository(db, viper.GetString("mongo.token")),
	}
}
