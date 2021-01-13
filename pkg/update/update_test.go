package update_test

import (
	"github.com/garethjevans/updatecenter/pkg/update"
	"testing"

	"github.com/garethjevans/updatecenter/pkg/api"

	"github.com/stretchr/testify/assert"
)

func TestCanParseUpdateCentre_SingleDependency(t *testing.T) {
	u := update.Updater{}

	http := &api.FakeHTTP{}
	client := api.NewClient(api.ReplaceTripper(http))
	u.SetClient(client)

	http.StubWithFixture(200, "golden.formatted.json")

	deps, err := u.LatestVersions([]string{"scm-api"})
	assert.NoError(t, err)

	assert.Equal(t, len(http.Requests), 1)
	assert.Equal(t, http.Requests[0].URL.Path, "/update-center.actual.json")
	assert.Equal(t, http.Requests[0].URL.Host, "updates.jenkins.io")
	assert.Equal(t, http.Requests[0].URL.RawQuery, "")

	t.Logf("got deps: %s", deps)

	assert.Equal(t, 2, len(deps))
	assertContains(t, deps, "scm-api")
	assertContains(t, deps, "structs")
	assertContainsWithVersion(t, deps, "structs", "1.20")
}

func TestCanParseUpdateCentre_MultipleDependency(t *testing.T) {
	u := update.Updater{}

	http := &api.FakeHTTP{}
	client := api.NewClient(api.ReplaceTripper(http))
	u.SetClient(client)

	http.StubWithFixture(200, "golden.formatted.json")

	deps, err := u.LatestVersions([]string{"scm-api", "job-dsl"})
	assert.NoError(t, err)

	assert.Equal(t, len(http.Requests), 1)
	assert.Equal(t, http.Requests[0].URL.Path, "/update-center.actual.json")
	assert.Equal(t, http.Requests[0].URL.Host, "updates.jenkins.io")
	assert.Equal(t, http.Requests[0].URL.RawQuery, "")

	t.Logf("got deps: %s", deps)

	assert.Equal(t, 4, len(deps))
	assertContains(t, deps, "scm-api")
	assertContains(t, deps, "structs")
	assertContainsWithVersion(t, deps, "structs", "1.20")
	assertContainsWithVersion(t, deps, "script-security", "1.54")
}

func assertContains(t *testing.T, deps []update.DepInfo, name string) {
	for _, a := range deps {
		if a.Name == name {
			return
		}
	}
	assert.Fail(t, "Was expecting to see %s", name)
}

func assertContainsWithVersion(t *testing.T, deps []update.DepInfo, name string, version string) {
	for _, a := range deps {
		if a.Name == name && a.Version == version{
			return
		}
	}
	assert.Fail(t, "Was expecting to see %s:%s", name, version)
}
