package handler

import (
	"github.com/pikachu0310/senirenol-server/core/internal/repository"
)

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// Common response
type ErrorResponse struct {
	Message string `json:"message"`
}
