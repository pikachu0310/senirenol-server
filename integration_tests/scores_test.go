package integrationtests

import (
	"encoding/json"
	"fmt"
	"testing"

	"gotest.tools/v3/assert"
)

func TestScores(t *testing.T) {
	// register user
	rec := doRequest(t, "POST", "/api/v1/users", "")
	assert.Equal(t, rec.Result().Status, `200 OK`)
	res := unmarshalResponse(t, rec)
	uid := res["id"].(string)

	// upsert chart
	rec = doRequest(t, "POST", "/api/v1/charts", `{"beatmap_id":"songY_future","song_name":"Song Y","difficulty":2}`)
	assert.Equal(t, rec.Result().Status, `200 OK`)

	// submit score
	body := fmt.Sprintf(`{"user_id":"%s","beatmap_id":"songY_future","score":987654,"max_combo":432,"perfect_critical_fast":100,"perfect_critical_late":120,"perfect_fast":50,"perfect_late":40,"good_fast":10,"good_late":5,"miss":3,"input":1}`, uid)
	rec = doRequest(t, "POST", "/api/v1/scores", body)
	assert.Equal(t, rec.Result().Status, `200 OK`)

	// ranking by chart
	rec = doRequest(t, "GET", "/api/v1/charts/ranking?beatmap_id=songY_future&limit=3", "")
	assert.Equal(t, rec.Result().Status, `200 OK`)
	// 単体でも配列で返る
	var arr []map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &arr)
	assert.Assert(t, len(arr) == 1)
	top := arr[0]["top"].([]any)
	assert.Assert(t, len(top) >= 1)
}
