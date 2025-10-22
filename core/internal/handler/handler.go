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

// Response DTOs
type (
	RegisterUserResponse struct {
		ID string `json:"id"`
	}

	UpdateUserNameResponse struct {
		Status string `json:"status"`
	}

	GetUserResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	UserStatsResponse struct {
		TotalPlays     int      `json:"total_plays"`
		DistinctCharts int      `json:"distinct_charts"`
		BestScore      *int     `json:"best_score"`
		AverageScore   *float64 `json:"average_score"`
	}

	UpsertChartResponse struct {
		Status string `json:"status"`
	}

	RankingEntryResponse struct {
		UserID     string `json:"user_id"`
		PlayerName string `json:"player_name"`
		Score      int    `json:"score"`
	}

	ChartRankingResponse struct {
		BeatmapID   string                 `json:"beatmap_id"`
		PlayerCount int                    `json:"player_count"`
		PlayCount   int                    `json:"play_count"`
		Top         []RankingEntryResponse `json:"top"`
	}

	SongPlaycountResponse struct {
		SongName  string `json:"song_name"`
		PlayCount int    `json:"play_count"`
	}
)
