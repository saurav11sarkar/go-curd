package notes

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repo struct {
	collection *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{db.Collection("notes")}
}

func (repo *Repo) Create(ctx context.Context, note Note) (Note, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := repo.collection.InsertOne(opCtx, note)
	if err != nil {
		return Note{}, err
	}
	return note, nil
}
