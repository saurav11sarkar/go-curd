package notes

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Handler struct {
	repo *Repo
}

func NewHandler(repo *Repo) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) CreateNote(ctx *gin.Context) {
	var req CreateNoteRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now().UTC()
	note := Note{
		ID:        bson.NewObjectID(),
		Title:     req.Title,
		Content:   req.Content,
		IsDeleted: false,
		Pinned:    req.Pinned,
		CreatedAt: now,
		UpdatedAt: now,
	}

	created, err := h.repo.Create(ctx.Request.Context(), note)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, created)
}

func (h *Handler) GetNotes(ctx *gin.Context) {
	notes, err := h.repo.FindAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, notes)
}

func (h *Handler) GetNote(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	note, err := h.repo.FindByID(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	ctx.JSON(http.StatusOK, note)
}

func (h *Handler) UpdateNote(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	var req UpdateNoteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := h.repo.Update(ctx.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid note id"})
		return
	}

	ctx.JSON(http.StatusOK, note)
}

func (h *Handler) DeleteNote(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	err := h.repo.Delete(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "note not found"})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
}
