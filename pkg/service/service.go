package service

import (
	"context"

	"github.com/qazaqpyn/bookCRUD/model"
	"github.com/qazaqpyn/bookCRUD/pkg/repository"
	audit "github.com/qazaqpyn/crud-audit-log/pkg/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Authorization interface {
	CreateUser(ctx context.Context, user model.User) error
	SignIn(ctx context.Context, inp model.LoginInput) (string, string, error)
	ParseToken(ctx context.Context, token string) (primitive.ObjectID, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
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

type AuditClient interface {
	SendLogRequest(ctx context.Context, req audit.LogItem) error
}

type RabbitMQServer interface {
	SendToQueue(msg model.Msg) error
}

func NewService(repos *repository.Repository, queue RabbitMQServer, audit AuditClient) *Service {
	return &Service{
		Authorization: NewAuthService(repos, queue),
		Book:          NewBookRepository(repos, audit),
	}
}
