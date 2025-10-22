package integrationtests

import (
	"encoding/json"
	"testing"

	"gotest.tools/v3/assert"
)

func TestCharts(t *testing.T) {
	// upsert chart
	rec := doRequest(t, "POST", "/api/v1/charts", `{"beatmap_id":"songX_parallel","song_name":"Song X","difficulty":4}`)
	t.Logf("charts upsert resp: %s", rec.Body.String())
	assert.Equal(t, rec.Result().Status, `200 OK`)

	// ranking all; JSONの基本構造を検証
	rec = doRequest(t, "GET", "/api/v1/charts/ranking?limit=5", "")
	assert.Equal(t, rec.Result().Status, `200 OK`)
	arr := []map[string]any{}
	_ = json.Unmarshal(rec.Body.Bytes(), &arr)
	assert.Assert(t, len(arr) >= 1)
	// 各要素に必要キーがあること
	for _, it := range arr {
		_, hasBeatmap := it["beatmap_id"]
		_, hasPlayCount := it["play_count"]
		assert.Assert(t, hasBeatmap)
		assert.Assert(t, hasPlayCount)
	}

	// songs playcount ranking
	rec = doRequest(t, "GET", "/api/v1/songs/playcount", "")
	assert.Equal(t, rec.Result().Status, `200 OK`)
}
