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

func TestScores_MultipleUsersAndPlays(t *testing.T) {
	// register users
	rec := doRequest(t, "POST", "/api/v1/users", "")
	assert.Equal(t, rec.Result().Status, `200 OK`)
	res := unmarshalResponse(t, rec)
	uid1 := res["id"].(string)

	rec = doRequest(t, "POST", "/api/v1/users", "")
	assert.Equal(t, rec.Result().Status, `200 OK`)
	res = unmarshalResponse(t, rec)
	uid2 := res["id"].(string)

	// update names for determinism
	rec = doRequest(t, "POST", "/api/v1/users/update", fmt.Sprintf(`{"user_id":"%s","user_name":"Alice"}`, uid1))
	assert.Equal(t, rec.Result().Status, `200 OK`)
	rec = doRequest(t, "POST", "/api/v1/users/update", fmt.Sprintf(`{"user_id":"%s","user_name":"Bob"}`, uid2))
	assert.Equal(t, rec.Result().Status, `200 OK`)

	// upsert chart
	rec = doRequest(t, "POST", "/api/v1/charts", `{"beatmap_id":"songM_future","song_name":"Song M","difficulty":2}`)
	assert.Equal(t, rec.Result().Status, `200 OK`)

	// submit many scores for uid1 (10 plays 980000..989000)
	for i := 0; i < 10; i++ {
		sc := 980000 + i*1000
		body := fmt.Sprintf(`{"user_id":"%s","beatmap_id":"songM_future","score":%d,"max_combo":300,"perfect_critical_fast":10,"perfect_critical_late":10,"perfect_fast":10,"perfect_late":10,"good_fast":1,"good_late":1,"miss":0,"input":0}`, uid1, sc)
		rec = doRequest(t, "POST", "/api/v1/scores", body)
		assert.Equal(t, rec.Result().Status, `200 OK`)
	}
	// submit some scores for uid2 (3 plays 970000..971000)
	for j := 0; j < 3; j++ {
		sc := 970000 + j*500
		body := fmt.Sprintf(`{"user_id":"%s","beatmap_id":"songM_future","score":%d,"max_combo":200,"perfect_critical_fast":8,"perfect_critical_late":8,"perfect_fast":8,"perfect_late":8,"good_fast":2,"good_late":2,"miss":1,"input":1}`, uid2, sc)
		rec = doRequest(t, "POST", "/api/v1/scores", body)
		assert.Equal(t, rec.Result().Status, `200 OK`)
	}

	// ranking by chart should aggregate correctly
	rec = doRequest(t, "GET", "/api/v1/charts/ranking?beatmap_id=songM_future&limit=10", "")
	assert.Equal(t, rec.Result().Status, `200 OK`)
	var arr []map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &arr)
	assert.Assert(t, len(arr) == 1)
	got := arr[0]
	assert.Equal(t, int(got["player_count"].(float64)), 2)
	assert.Equal(t, int(got["play_count"].(float64)), 13)
	// top order: Alice first with 989000, then Bob with 971000
	top := got["top"].([]any)
	assert.Assert(t, len(top) == 2)
	top0 := top[0].(map[string]any)
	top1 := top[1].(map[string]any)
	assert.Equal(t, top0["player_name"].(string), "Alice")
	assert.Equal(t, int(top0["score"].(float64)), 989000)
	assert.Equal(t, top1["player_name"].(string), "Bob")
	assert.Equal(t, int(top1["score"].(float64)), 971000)

	// song playcount ranking should include Song M with 13 plays
	rec = doRequest(t, "GET", "/api/v1/songs/playcount", "")
	assert.Equal(t, rec.Result().Status, `200 OK`)
	var songArr []map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &songArr)
	found := false
	for _, it := range songArr {
		nameVal, ok := it["song_name"]
		if !ok || nameVal == nil {
			continue
		}
		name, ok := nameVal.(string)
		if !ok {
			continue
		}
		if name == "Song M" {
			pcVal, ok := it["play_count"].(float64)
			assert.Assert(t, ok)
			assert.Equal(t, int(pcVal), 13)
			found = true
		}
	}
	assert.Assert(t, found)
}
