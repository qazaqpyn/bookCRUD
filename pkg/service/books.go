package service

import (
	"context"

	"github.com/qazaqpyn/bookCRUD/model"
	"github.com/qazaqpyn/bookCRUD/pkg/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookService struct {
	repo *repository.Repository
}

func NewBookRepository(repo *repository.Repository) *BookService {
	return &BookService{
		repo: repo,
	}
}

func (b *BookService) Create(ctx context.Context, book model.Book) error {
	book.Id = primitive.NewObjectID()
	return b.repo.CreateBook(ctx, book)
}

func (b *BookService) GetById(ctx context.Context, id primitive.ObjectID) (*model.Book, error) {
	return b.repo.GetBookById(ctx, id)
}

func (b *BookService) GetAll(ctx context.Context) ([]*model.Book, error) {
	return b.repo.GetAllBook(ctx)
}

func (b *BookService) Update(ctx context.Context, id primitive.ObjectID, book *model.Book) error {
	return b.repo.UpdateBook(ctx, id, book)
}

func (b *BookService) Delete(ctx context.Context, id primitive.ObjectID) error {
	return b.repo.DeleteBook(ctx, id)
}
