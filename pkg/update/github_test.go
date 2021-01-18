package update_test

import (
	"testing"

	"github.com/garethjevans/uc/pkg/api"
	"github.com/garethjevans/uc/pkg/update"
	"github.com/stretchr/testify/assert"
)

func TestGetLatestLTSRelease(t *testing.T) {
	g := update.GitHub{}

	http := &api.FakeHTTP{}
	client := api.NewClient(api.ReplaceTripper(http))
	g.SetClient(client)

	http.StubWithFixture(200, "releases.json")

	version, err := g.GetLatestLTSRelease()
	assert.NoError(t, err)
	assert.Equal(t, "2.263.2", version)
}
