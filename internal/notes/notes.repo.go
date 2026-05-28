package notes

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repo struct {
	collection *mongo.Collection
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{collection: db.Collection("notes")}
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

func (repo *Repo) FindAll(ctx context.Context) ([]Note, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := repo.collection.Find(opCtx, bson.M{"isDeleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(opCtx)

	var notes []Note
	if err := cursor.All(opCtx, &notes); err != nil {
		return nil, err
	}

	if notes == nil {
		notes = []Note{}
	}

	return notes, nil
}

func (repo *Repo) FindByID(ctx context.Context, id string) (Note, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return Note{}, err
	}

	var note Note
	err = repo.collection.FindOne(opCtx, bson.M{
		"_id": objID,
	}).Decode(&note)
	if err != nil {
		return Note{}, err
	}

	return note, nil
}

func (repo *Repo) Update(ctx context.Context, id string, req UpdateNoteRequest) (Note, error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return Note{}, err
	}
	allowed := map[string]bool{
		"title":     true,
		"content":   true,
		"isDeleted": true,
		"pinned":    true,
	}
	set := bson.M{}
	for key, value := range req {
		if allowed[key] {
			set[key] = value
		}
	}

	set["updatedAt"] = time.Now().UTC()

	var note Note
	err = repo.collection.FindOneAndUpdate(
		opCtx,
		bson.M{"_id": objID},
		bson.M{"$set": set},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&note)
	if err != nil {
		return Note{}, err
	}

	return note, nil
}

func (repo *Repo) Delete(ctx context.Context, id string) error {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := repo.collection.DeleteOne(opCtx, bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
