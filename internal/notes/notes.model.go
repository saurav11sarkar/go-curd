package notes

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Note struct {
	ID        bson.ObjectID `json:"id" bson:"_id"`
	Title     string        `json:"title" bson:"title"`
	Content   string        `json:"content" bson:"content"`
	IsDeleted bool          `json:"isDeleted" bson:"isDeleted"`
	Pinned    bool          `json:"pinned" bson:"pinned"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
}

type CreateNoteRequest struct {
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	IsDeleted bool   `json:"isDeleted"`
	Pinned    bool   `json:"pinned"`
}

type UpdateNoteRequest map[string]any
