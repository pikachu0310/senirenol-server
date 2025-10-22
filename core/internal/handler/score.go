package handler

import (
	"net/http"

	vd "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
	"github.com/pikachu0310/senirenol-server/core/internal/repository"
)

type SubmitScoreRequest struct {
	UserID              string `json:"user_id"`
	BeatmapID           string `json:"beatmap_id"`
	Score               int    `json:"score"`
	MaxCombo            int    `json:"max_combo"`
	PerfectCriticalFast int    `json:"perfect_critical_fast"`
	PerfectCriticalLate int    `json:"perfect_critical_late"`
	PerfectFast         int    `json:"perfect_fast"`
	PerfectLate         int    `json:"perfect_late"`
	GoodFast            int    `json:"good_fast"`
	GoodLate            int    `json:"good_late"`
	Miss                int    `json:"miss"`
	Input               uint8  `json:"input"`
}

type SubmitScoreResponse struct {
	ID int64 `json:"id"`
}

// SubmitScore godoc
// @Summary スコア登録
// @Description プレイ結果を登録します
// @Tags scores
// @Accept json
// @Produce json
// @Param score body SubmitScoreRequest true "スコア情報"
// @Success 200 {object} SubmitScoreResponse "登録されたスコアのID"
// @Failure 400 {object} ErrorResponse
// @Router /scores [post]
func (h *Handler) SubmitScore(c echo.Context) error {
	var req SubmitScoreRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").SetInternal(err)
	}
	if err := vd.ValidateStruct(&req,
		vd.Field(&req.UserID, vd.Required),
		vd.Field(&req.BeatmapID, vd.Required),
	); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	id, err := h.repo.InsertScore(c.Request().Context(), repository.InsertScoreParams{
		UserID:              req.UserID,
		BeatmapID:           req.BeatmapID,
		Score:               req.Score,
		MaxCombo:            req.MaxCombo,
		PerfectCriticalFast: req.PerfectCriticalFast,
		PerfectCriticalLate: req.PerfectCriticalLate,
		PerfectFast:         req.PerfectFast,
		PerfectLate:         req.PerfectLate,
		GoodFast:            req.GoodFast,
		GoodLate:            req.GoodLate,
		Miss:                req.Miss,
		Input:               repository.InputType(req.Input),
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return c.JSON(http.StatusOK, SubmitScoreResponse{ID: id})
}
