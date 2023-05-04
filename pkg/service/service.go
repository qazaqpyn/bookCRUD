package service

import (
	"context"

	"github.com/qazaqpyn/bookCRUD/model"
	"github.com/qazaqpyn/bookCRUD/pkg/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Authorization interface {
	CreateUser(ctx context.Context, user model.User) error
	GenerateToken(ctx context.Context, email, password string) (string, error)
	ParseToken(ctx context.Context, token string) (primitive.ObjectID, error)
}

type Book interface {
	Create(ctx context.Context, book model.Book) error
	GetById(ctx context.Context, id primitive.ObjectID) (*model.Book, error)
	GetAll(ctx context.Context) ([]*model.Book, error)
	Update(ctx context.Context, id primitive.ObjectID, book *model.Book) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type Service struct {
	Authorization
	Book
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos),
		Book:          NewBookRepository(repos),
	}
}
