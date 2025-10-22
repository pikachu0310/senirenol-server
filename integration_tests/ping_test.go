package integrationtests

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestPing(t *testing.T) {
	t.Run("pong", func(t *testing.T) {
		t.Parallel()
		rec := doRequest(t, "GET", "/api/v1/ping", "")
		assert.Equal(t, rec.Result().Status, `200 OK`)
		assert.Equal(t, rec.Body.String(), "pong")
	})
}
