-- +goose Up

-- users: Senirenol Bloom players
CREATE TABLE IF NOT EXISTS users (
	id VARCHAR(36) NOT NULL,
	name VARCHAR(255) NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
);

-- charts: beatmaps for songs and difficulties
-- difficulty enum: 0=Past,1=Present,2=Future,3=Lycoris,4=Parallel
CREATE TABLE IF NOT EXISTS charts (
	beatmap_id VARCHAR(128) NOT NULL,
	song_name VARCHAR(255) NOT NULL,
	difficulty TINYINT NOT NULL,
	parallel_string VARCHAR(16) NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (beatmap_id)
);

-- scores: individual play results
-- input enum: 0=keyboard,1=button
CREATE TABLE IF NOT EXISTS scores (
	id BIGINT NOT NULL AUTO_INCREMENT,
	user_id VARCHAR(36) NOT NULL,
	beatmap_id VARCHAR(128) NOT NULL,
	score INT NOT NULL,
	max_combo INT NOT NULL,
	perfect_critical_fast INT NOT NULL,
	perfect_critical_late INT NOT NULL,
	perfect_fast INT NOT NULL,
	perfect_late INT NOT NULL,
	good_fast INT NOT NULL,
	good_late INT NOT NULL,
	miss INT NOT NULL,
	input TINYINT NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	KEY idx_scores_beatmap (beatmap_id),
	KEY idx_scores_user (user_id),
	KEY idx_scores_beatmap_score (beatmap_id, score DESC),
	CONSTRAINT fk_scores_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
	CONSTRAINT fk_scores_chart FOREIGN KEY (beatmap_id) REFERENCES charts(beatmap_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS scores;
DROP TABLE IF EXISTS charts;
DROP TABLE IF EXISTS users;
