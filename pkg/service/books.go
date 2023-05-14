package service

import (
	"context"
	"errors"
	"time"

	"github.com/qazaqpyn/bookCRUD/model"
	"github.com/qazaqpyn/bookCRUD/pkg/logging"
	"github.com/qazaqpyn/bookCRUD/pkg/repository"
	audit "github.com/qazaqpyn/crud-audit-log/pkg/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookService struct {
	repo        *repository.Repository
	auditClient AuditClient
}

func NewBookRepository(repo *repository.Repository, auditClient AuditClient) *BookService {
	return &BookService{
		repo:        repo,
		auditClient: auditClient,
	}
}

func (b *BookService) Create(ctx context.Context, book model.Book) error {
	book.Id = primitive.NewObjectID()
	err := b.repo.CreateBook(ctx, book)
	if err != nil {
		return err
	}

	userID, ok := ctx.Value("userId").(string)
	if !ok {
		return errors.New("user not authorized")
	}

	if err := b.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_CREATE,
		Entity:    audit.ENTITY_BOOK,
		EntityID:  userID,
		Timestamp: time.Now(),
	}); err != nil {
		logging.LogError("Book.Create", err)
	}

	return nil
}

func (b *BookService) GetById(ctx context.Context, id primitive.ObjectID) (*model.Book, error) {
	book, err := b.repo.GetBookById(ctx, id)
	if err != nil {
		return nil, err
	}

	userID, ok := ctx.Value("userId").(string)
	if !ok {
		return nil, errors.New("user not authorized")
	}

	if err := b.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_GET,
		Entity:    audit.ENTITY_BOOK,
		EntityID:  userID,
		Timestamp: time.Now(),
	}); err != nil {
		logging.LogError("Book.GetByID", err)
	}

	return book, nil
}

func (b *BookService) GetAll(ctx context.Context) ([]*model.Book, error) {
	books, err := b.repo.GetAllBook(ctx)
	if err != nil {
		return nil, err
	}

	userID, ok := ctx.Value("userId").(string)
	if !ok {
		return nil, errors.New("user not authorized")
	}

	if err := b.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_GET,
		Entity:    audit.ENTITY_BOOK,
		EntityID:  userID,
		Timestamp: time.Now(),
	}); err != nil {
		logging.LogError("Book.GetAll", err)
	}

	return books, nil
}

func (b *BookService) Update(ctx context.Context, id primitive.ObjectID, book *model.Book) error {
	err := b.repo.UpdateBook(ctx, id, book)
	if err != nil {
		return err
	}

	userID, ok := ctx.Value("userId").(string)
	if !ok {
		return errors.New("user not authorized")
	}

	if err := b.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_UPDATE,
		Entity:    audit.ENTITY_BOOK,
		EntityID:  userID,
		Timestamp: time.Now(),
	}); err != nil {
		logging.LogError("Book.Update", err)
	}

	return nil
}

func (b *BookService) Delete(ctx context.Context, id primitive.ObjectID) error {
	err := b.repo.DeleteBook(ctx, id)
	if err != nil {
		return err
	}

	userID, ok := ctx.Value("userId").(string)
	if !ok {
		return errors.New("user not authorized")
	}

	if err := b.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ACTION_DELETE,
		Entity:    audit.ENTITY_BOOK,
		EntityID:  userID,
		Timestamp: time.Now(),
	}); err != nil {
		logging.LogError("Book.Delete", err)
	}

	return nil
}
