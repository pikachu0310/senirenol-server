package integrationtests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"gotest.tools/v3/assert"
)

func TestUsers(t *testing.T) {
	// register
	rec := doRequest(t, "POST", "/api/v1/users", "")
	t.Logf("register resp: %s", rec.Body.String())
	assert.Equal(t, rec.Result().Status, `200 OK`)
	res := unmarshalResponse(t, rec)
	uid := res["id"].(string)
	_ = uuid.MustParse(uid)

	// update name
	rec = doRequest(t, "POST", "/api/v1/users/update", fmt.Sprintf(`{"user_id":"%s","user_name":"Bob"}`, uid))
	t.Logf("update resp: %s", rec.Body.String())
	assert.Equal(t, rec.Result().Status, `200 OK`)

	// get user
	rec = doRequest(t, "GET", "/api/v1/users/"+uid, "")
	t.Logf("get user resp: %s", rec.Body.String())
	assert.Equal(t, rec.Result().Status, `200 OK`)

	// stats (no scores yet also OK)
	rec = doRequest(t, "GET", "/api/v1/users/"+uid+"/stats", "")
	t.Logf("stats resp: %s", rec.Body.String())
	assert.Equal(t, rec.Result().Status, `200 OK`)
}
