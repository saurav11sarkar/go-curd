package notes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func RegisterRoutes(r *gin.Engine, db *mongo.Database) {
	repo := NewRepo(db)
	h := NewHandler(repo)

	notesGroup := r.Group("/notes")
	{
		notesGroup.POST("/", h.CreateNote)
		notesGroup.GET("/", h.GetNotes)
		notesGroup.GET("/:id", h.GetNote)
		notesGroup.PUT("/:id", h.UpdateNote)
		notesGroup.DELETE("/:id", h.DeleteNote)
	}
}
