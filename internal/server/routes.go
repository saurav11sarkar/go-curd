package server

import (
	"curd/internal/notes"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewRoutes(database *mongo.Database) *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong server is connect",
		})
	})

	notes.RegisterRoutes(router, database)
	return router
}
