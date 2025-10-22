// NOTE: go test -update でスナップショット更新可

package integrationtests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"gotest.tools/v3/assert"
)

func TestSenirenolFlow(t *testing.T) {
	// 1) ユーザー登録
	rec := doRequest(t, "POST", "/api/v1/users", "")
	assert.Equal(t, rec.Result().Status, `200 OK`)
	res := unmarshalResponse(t, rec)
	uid := res["id"].(string)
	_ = uuid.MustParse(uid)

	// 2) ユーザー名更新
	rec = doRequest(t, "POST", "/api/v1/users/update", fmt.Sprintf(`{"user_id":"%s","user_name":"Alice"}`, uid))
	assert.Equal(t, rec.Result().Status, `200 OK`)

	// 3) 譜面登録
	rec = doRequest(t, "POST", "/api/v1/charts", `{"beatmap_id":"song1_future","song_name":"Song 1","difficulty":2}`)
	assert.Equal(t, rec.Result().Status, `200 OK`)

	// 4) スコア投稿
	rec = doRequest(t, "POST", "/api/v1/scores", fmt.Sprintf(`{
        "user_id":"%s",
        "beatmap_id":"song1_future",
        "score":990000,
        "max_combo":500,
        "perfect_critical_fast":200,
        "perfect_critical_late":210,
        "perfect_fast":50,
        "perfect_late":30,
        "good_fast":5,
        "good_late":3,
        "miss":2,
        "input":0
    }`, uid))
	assert.Equal(t, rec.Result().Status, `200 OK`)

	// 5) 譜面ランキング取得（個別）
	rec = doRequest(t, "GET", "/api/v1/charts/ranking?beatmap_id=song1_future&limit=10", "")
	assert.Equal(t, rec.Result().Status, `200 OK`)

	// 6) 楽曲プレイ回数ランキング
	rec = doRequest(t, "GET", "/api/v1/songs/playcount", "")
	assert.Equal(t, rec.Result().Status, `200 OK`)

	// 7) ユーザー統計
	rec = doRequest(t, "GET", "/api/v1/users/"+uid+"/stats", "")
	assert.Equal(t, rec.Result().Status, `200 OK`)
}
