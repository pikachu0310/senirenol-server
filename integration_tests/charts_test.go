package integrationtests

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestCharts(t *testing.T) {
	// upsert chart
	rec := doRequest(t, "POST", "/api/v1/charts", `{"beatmap_id":"songX_parallel","song_name":"Song X","difficulty":4}`)
	t.Logf("charts upsert resp: %s", rec.Body.String())
	assert.Equal(t, rec.Result().Status, `200 OK`)

	// ranking all
	rec = doRequest(t, "GET", "/api/v1/charts/ranking?limit=5", "")
	t.Logf("charts ranking resp: %s", rec.Body.String())
	assert.Equal(t, rec.Result().Status, `200 OK`)

	// songs playcount ranking
	rec = doRequest(t, "GET", "/api/v1/songs/playcount", "")
	assert.Equal(t, rec.Result().Status, `200 OK`)
}
