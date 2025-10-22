package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Chart struct {
	BeatmapID      string     `db:"beatmap_id"`
	SongName       string     `db:"song_name"`
	Difficulty     Difficulty `db:"difficulty"`
	ParallelString *string    `db:"parallel_string"`
	CreatedAt      time.Time  `db:"created_at"`
}

type UpsertChartParams struct {
	BeatmapID      string
	SongName       string
	Difficulty     Difficulty
	ParallelString *string
}

func (r *Repository) UpsertChart(ctx context.Context, p UpsertChartParams) error {
	// MySQL upsert
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO charts (beatmap_id, song_name, difficulty, parallel_string)
        VALUES (?, ?, ?, ?)
        ON DUPLICATE KEY UPDATE song_name=VALUES(song_name), difficulty=VALUES(difficulty), parallel_string=VALUES(parallel_string)
    `, p.BeatmapID, p.SongName, p.Difficulty, p.ParallelString)
	if err != nil {
		return fmt.Errorf("upsert chart: %w", err)
	}
	return nil
}

func (r *Repository) GetChart(ctx context.Context, beatmapID string) (*Chart, error) {
	var c Chart
	if err := r.db.GetContext(ctx, &c, `SELECT * FROM charts WHERE beatmap_id=?`, beatmapID); err != nil {
		return nil, fmt.Errorf("get chart: %w", err)
	}
	return &c, nil
}

func (r *Repository) GetCharts(ctx context.Context) ([]*Chart, error) {
	var cs []*Chart
	if err := r.db.SelectContext(ctx, &cs, `SELECT * FROM charts ORDER BY song_name, difficulty, beatmap_id`); err != nil {
		return nil, fmt.Errorf("get charts: %w", err)
	}
	return cs, nil
}

type ChartStats struct {
	BeatmapID   string   `db:"beatmap_id"`
	PlayCount   int      `db:"play_count"`
	PlayerCount int      `db:"player_count"`
	AvgScore    *float64 `db:"avg_score"`
	BestScore   *int     `db:"best_score"`
}

func (r *Repository) GetChartStats(ctx context.Context, beatmapID string) (*ChartStats, error) {
	var s ChartStats
	// Aggregates with NULL-safe handling
	if err := r.db.GetContext(ctx, &s, `
        SELECT ? AS beatmap_id,
               COUNT(*) AS play_count,
               COUNT(DISTINCT user_id) AS player_count,
               AVG(score) AS avg_score,
               MAX(score) AS best_score
        FROM scores WHERE beatmap_id = ?
    `, beatmapID, beatmapID); err != nil {
		if err == sql.ErrNoRows {
			return &ChartStats{BeatmapID: beatmapID, PlayCount: 0, PlayerCount: 0, AvgScore: nil, BestScore: nil}, nil
		}
		return nil, fmt.Errorf("chart stats: %w", err)
	}
	return &s, nil
}

type SongPlayCount struct {
	SongName  string `db:"song_name"`
	PlayCount int    `db:"play_count"`
}

func (r *Repository) GetSongPlaycountRanking(ctx context.Context) ([]*SongPlayCount, error) {
	var rs []*SongPlayCount
	if err := r.db.SelectContext(ctx, &rs, `
        SELECT c.song_name, COUNT(*) AS play_count
        FROM scores s
        JOIN charts c ON c.beatmap_id = s.beatmap_id
        GROUP BY c.song_name
        ORDER BY play_count DESC, c.song_name ASC
    `); err != nil {
		return nil, fmt.Errorf("song playcount ranking: %w", err)
	}
	return rs, nil
}
