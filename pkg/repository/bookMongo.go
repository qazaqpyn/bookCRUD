package repository

import (
	"context"

	"github.com/qazaqpyn/bookCRUD/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookRepository struct {
	db *mongo.Collection
}

func NewBookRepository(db *mongo.Database, collection string) *BookRepository {
	return &BookRepository{
		db: db.Collection(collection),
	}
}

func (r BookRepository) CreateBook(ctx context.Context, book model.Book) error {
	_, err := r.db.InsertOne(ctx, book)
	if err != nil {
		return err
	}

	return nil
}

func (r BookRepository) GetBookById(ctx context.Context, id primitive.ObjectID) (*model.Book, error) {
	book := new(model.Book)

	err := r.db.FindOne(ctx, bson.M{
		"id": id,
	}).Decode(book)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (r BookRepository) GetAllBook(ctx context.Context) ([]*model.Book, error) {
	cur, err := r.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	out := make([]*model.Book, 0)

	for cur.Next(ctx) {
		book := new(model.Book)
		err := cur.Decode(book)
		if err != nil {
			return nil, err
		}

		out = append(out, book)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func (r BookRepository) DeleteBook(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.db.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (r BookRepository) UpdateBook(ctx context.Context, id primitive.ObjectID, book *model.Book) error {
	_, err := r.db.UpdateByID(ctx, bson.M{"id": id}, book)
	return err
}
