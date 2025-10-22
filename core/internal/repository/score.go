package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

type InsertScoreParams struct {
	UserID              string
	BeatmapID           string
	Score               int
	MaxCombo            int
	PerfectCriticalFast int
	PerfectCriticalLate int
	PerfectFast         int
	PerfectLate         int
	GoodFast            int
	GoodLate            int
	Miss                int
	Input               InputType
}

type ScoreRow struct {
	ID                  int64     `db:"id"`
	UserID              string    `db:"user_id"`
	BeatmapID           string    `db:"beatmap_id"`
	Score               int       `db:"score"`
	MaxCombo            int       `db:"max_combo"`
	PerfectCriticalFast int       `db:"perfect_critical_fast"`
	PerfectCriticalLate int       `db:"perfect_critical_late"`
	PerfectFast         int       `db:"perfect_fast"`
	PerfectLate         int       `db:"perfect_late"`
	GoodFast            int       `db:"good_fast"`
	GoodLate            int       `db:"good_late"`
	Miss                int       `db:"miss"`
	Input               InputType `db:"input"`
	CreatedAt           time.Time `db:"created_at"`
}

func (r *Repository) InsertScore(ctx context.Context, p InsertScoreParams) (int64, error) {
	res, err := r.db.ExecContext(ctx, `
        INSERT INTO scores (
            user_id, beatmap_id, score, max_combo,
            perfect_critical_fast, perfect_critical_late,
            perfect_fast, perfect_late,
            good_fast, good_late,
            miss, input
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `, p.UserID, p.BeatmapID, p.Score, p.MaxCombo, p.PerfectCriticalFast, p.PerfectCriticalLate, p.PerfectFast, p.PerfectLate, p.GoodFast, p.GoodLate, p.Miss, p.Input)
	if err != nil {
		return 0, fmt.Errorf("insert score: %w", err)
	}
	id, _ := res.LastInsertId()
	return id, nil
}

type RankingEntry struct {
	UserID string `db:"user_id" json:"user_id"`
	Name   string `db:"name" json:"player_name"`
	Score  int    `db:"best_score" json:"score"`
}

type ChartRanking struct {
	BeatmapID   string         `json:"beatmap_id"`
	PlayerCount int            `json:"player_count"`
	PlayCount   int            `json:"play_count"`
	Top         []RankingEntry `json:"top"`
}

func (r *Repository) GetChartRanking(ctx context.Context, beatmapID string, limit int) (*ChartRanking, error) {
	// player_count, play_count
	var counts struct {
		PlayerCount int `db:"player_count"`
		PlayCount   int `db:"play_count"`
	}
	if err := r.db.GetContext(ctx, &counts, `
        SELECT COUNT(DISTINCT user_id) AS player_count, COUNT(*) AS play_count
        FROM scores WHERE beatmap_id = ?
    `, beatmapID); err != nil {
		return nil, fmt.Errorf("count ranking: %w", err)
	}

	var top []RankingEntry
	query := `
        SELECT s.user_id, u.name, MAX(s.score) AS best_score
        FROM scores s
        JOIN users u ON u.id = s.user_id
        WHERE s.beatmap_id = ?
        GROUP BY s.user_id, u.name
        ORDER BY best_score DESC, u.name ASC
        LIMIT ` + strconv.Itoa(limit)
	if err := r.db.SelectContext(ctx, &top, query, beatmapID); err != nil {
		return nil, fmt.Errorf("top ranking: %w", err)
	}

	return &ChartRanking{
		BeatmapID:   beatmapID,
		PlayerCount: counts.PlayerCount,
		PlayCount:   counts.PlayCount,
		Top:         top,
	}, nil
}

func (r *Repository) GetAllChartsRankings(ctx context.Context, limit int) ([]*ChartRanking, error) {
	charts, err := r.GetCharts(ctx)
	if err != nil {
		// 取得に失敗した場合は空配列を返す
		return []*ChartRanking{}, nil
	}
	res := make([]*ChartRanking, 0, len(charts))
	for _, c := range charts {
		cr, err := r.GetChartRanking(ctx, c.BeatmapID, limit)
		if err != nil {
			// 個別の失敗はスキップ
			continue
		}
		res = append(res, cr)
	}
	return res, nil
}

type UserStats struct {
	TotalPlays     int      `db:"total_plays"`
	DistinctCharts int      `db:"distinct_charts"`
	BestScore      *int     `db:"best_score"`
	AverageScore   *float64 `db:"avg_score"`
}

func (r *Repository) GetUserStats(ctx context.Context, userID string) (*UserStats, error) {
	var s UserStats
	if err := r.db.GetContext(ctx, &s, `
        SELECT COUNT(*) AS total_plays,
               COUNT(DISTINCT beatmap_id) AS distinct_charts,
               MAX(score) AS best_score,
               AVG(score) AS avg_score
        FROM scores WHERE user_id = ?
    `, userID); err != nil {
		return nil, fmt.Errorf("user stats: %w", err)
	}
	return &s, nil
}
