package handler

import (
	"net/http"
	"strconv"

	vd "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
	"github.com/pikachu0310/senirenol-server/core/internal/repository"
)

type UpsertChartRequest struct {
	BeatmapID      string  `json:"beatmap_id"`
	SongName       string  `json:"song_name"`
	Difficulty     uint8   `json:"difficulty"`
	ParallelString *string `json:"parallel_string"`
}

// UpsertChart godoc
// @Summary 譜面登録/更新
// @Tags charts
// @Accept json
// @Produce json
// @Param chart body UpsertChartRequest true "譜面情報"
// @Success 200 {object} UpsertChartResponse "登録結果"
// @Failure 400 {object} ErrorResponse
// @Router /charts [post]
func (h *Handler) UpsertChart(c echo.Context) error {
	req := new(UpsertChartRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body").SetInternal(err)
	}
	if err := vd.ValidateStruct(
		req,
		vd.Field(&req.BeatmapID, vd.Required),
		vd.Field(&req.SongName, vd.Required),
		vd.Field(&req.Difficulty, vd.Required),
	); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := h.repo.UpsertChart(c.Request().Context(), repository.UpsertChartParams{
		BeatmapID:      req.BeatmapID,
		SongName:       req.SongName,
		Difficulty:     int(req.Difficulty),
		ParallelString: req.ParallelString,
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	return c.JSON(http.StatusOK, UpsertChartResponse{Status: "ok"})
}

// GetChartRanking godoc
// @Summary 譜面ランキング
// @Description beatmap_idを指定したランキング、未指定時は全譜面のランキング
// @Tags charts
// @Produce json
// @Param beatmap_id query string false "譜面ID"
// @Param limit query int false "ランキング上限"
// @Success 200 {array} ChartRankingResponse "ランキング配列"
// @Router /charts/ranking [get]
func (h *Handler) GetChartRanking(c echo.Context) error {
	beatmapID := c.QueryParam("beatmap_id")
	limitStr := c.QueryParam("limit")
	limit := 10
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
			limit = v
		}
	}
	if beatmapID != "" {
		r, err := h.repo.GetChartRanking(c.Request().Context(), beatmapID, limit)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
		}
		// 単体でも配列で返す
		return c.JSON(http.StatusOK, []ChartRankingResponse{{
			BeatmapID:   r.BeatmapID,
			PlayerCount: r.PlayerCount,
			PlayCount:   r.PlayCount,
			Top:         toRankingEntryResponse(r.Top),
		}})
	}
	rs, err := h.repo.GetAllChartsRankings(c.Request().Context(), limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	out := make([]ChartRankingResponse, 0, len(rs))
	for _, r := range rs {
		out = append(out, ChartRankingResponse{
			BeatmapID:   r.BeatmapID,
			PlayerCount: r.PlayerCount,
			PlayCount:   r.PlayCount,
			Top:         toRankingEntryResponse(r.Top),
		})
	}
	return c.JSON(http.StatusOK, out)
}

// GetSongPlaycountRanking godoc
// @Summary 楽曲プレイ回数ランキング
// @Tags charts
// @Produce json
// @Success 200 {array} map[string]any
// @Router /songs/playcount [get]
func (h *Handler) GetSongPlaycountRanking(c echo.Context) error {
	rs, err := h.repo.GetSongPlaycountRanking(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	out := make([]SongPlaycountResponse, len(rs))
	for i, r := range rs {
		out[i] = SongPlaycountResponse{SongName: r.SongName, PlayCount: r.PlayCount}
	}
	return c.JSON(http.StatusOK, out)
}

func toRankingEntryResponse(in []repository.RankingEntry) []RankingEntryResponse {
	if in == nil {
		return nil
	}
	out := make([]RankingEntryResponse, len(in))
	for i, e := range in {
		out[i] = RankingEntryResponse{
			UserID:     e.UserID,
			PlayerName: e.Name,
			Score:      e.Score,
		}
	}
	return out
}
