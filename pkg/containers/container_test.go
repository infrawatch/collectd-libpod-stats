package containers

import (
	"encoding/json"
	"testing"

	"github.com/pleimer/collectd-libpod-stats/pkg/assert"
)

func TestListCreation(t *testing.T) {
	c := Container{
		Names:   []string{"qdr, qdr_metric"},
		ID:      "testing-id",
		Image:   "long-string-of-characters",
		Created: "some-utc-timestamp",
	}
	cList := []*Container{&c}

	cListJSON, err := json.Marshal(cList)
	assert.Ok(t, err)

	t.Run("normal operation", func(t *testing.T) {
		mUT, err := NewMapFromJSON(cListJSON)
		assert.Ok(t, err)

		std := List{
			c.ID: &c,
		}

		assert.Equals(t, std, mUT)
	})

	t.Run("bad json data", func(t *testing.T) {
		badJSON := cListJSON
		badJSON = append(badJSON, []byte{130}...)

		_, err := NewMapFromJSON(badJSON)
		assert.Assert(t, (err != nil), "expected a json error")
	})

	t.Run("wrong format", func(t *testing.T) {
		cJSON, err := json.Marshal(c)
		assert.Ok(t, err)

		_, err = NewMapFromJSON(cJSON)
		assert.Assert(t, (err != nil), "expected a json error")
	})

	t.Run("container ID not created", func(t *testing.T) {
		badID := cList
		badID[0].ID = ""
		badIDJSON, err := json.Marshal(badID)
		assert.Ok(t, err)

		_, err = NewMapFromJSON(badIDJSON)
		assert.Assert(t, (err != nil), "expected invalid ID error")
	})
}
