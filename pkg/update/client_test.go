package update_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/garethjevans/uc/pkg/update"

	"github.com/garethjevans/uc/pkg/api"

	"github.com/stretchr/testify/assert"
)

var (
	versionSet = `ansicolor:0.7.3
antisamy-markup-formatter:2.1
authentication-tokens:1.4
azure-container-agents:1.2.1
azure-vm-agents:1.5.2
basic-branch-build-strategies:1.3.2
blueocean:1.24.3
blueocean-autofavorite:1.2.4
blueocean-commons:1.24.3
blueocean-config:1.24.3
blueocean-core-js:1.24.3
blueocean-dashboard:1.24.3
blueocean-display-url:2.4.0
blueocean-events:1.24.3
blueocean-git-pipeline:1.24.3
blueocean-github-pipeline:1.24.3
blueocean-i18n:1.24.3
blueocean-jira:1.24.3
blueocean-jwt:1.24.3
blueocean-personalization:1.24.3
blueocean-pipeline-api-impl:1.24.3
blueocean-pipeline-editor:1.24.3
blueocean-pipeline-scm-api:1.24.3
blueocean-rest:1.24.3
blueocean-rest-impl:1.24.3
blueocean-web:1.24.3
branch-api:2.6.3
build-name-setter:2.1.0
cloud-stats:0.26
config-file-provider:3.7.0
configuration-as-code:1.46
credentials:2.3.14
credentials-binding:1.24
datadog:2.6.0
ec2:1.56
embeddable-build-status:2.0.3
extended-read-permission:3.2
git:4.5.2
git-client:3.6.0
github:1.32.0
github-api:1.117
github-branch-source:2.9.3
github-checks:1.0.8
github-label-filter:1.0.0
groovy:2.3
inline-pipeline:1.0.1
javadoc:1.6
jira:3.1.3
job-dsl:1.77
junit:1.48
kubernetes:1.28.5
kubernetes-credentials:0.7.0
kubernetes-credentials-provider:0.15
ldap:2.2
lockable-resources:2.10
matrix-auth:2.6.4
matrix-project:1.18
metrics:4.0.2.7
pipeline-build-step:2.13
pipeline-github:2.7
pipeline-graph-analysis:1.10
pipeline-input-step:2.12
pipeline-milestone-step:1.3.1
pipeline-model-api:1.7.2
pipeline-model-definition:1.7.2
pipeline-model-extensions:1.7.2
pipeline-rest-api:2.19
pipeline-stage-step:2.5
pipeline-stage-tags-metadata:1.7.2
pipeline-stage-view:2.19
pipeline-utility-steps:2.6.1
plain-credentials:1.7
prometheus:2.0.8
scm-api:2.6.4
scm-filter-branch-pr:0.5.1
script-security:1.75
ssh-agent:1.20
ssh-credentials:1.18.1
support-core:2.72
timestamper:1.11.8
token-macro:2.14
toolenv:1.2
variant:1.4
warnings-ng:8.6.0
workflow-aggregator:2.6
workflow-api:2.40
workflow-basic-steps:2.23
workflow-cps:2.87
workflow-cps-global-lib:2.17
workflow-durable-task-step:2.37
workflow-job:2.40
workflow-multibranch:2.22
workflow-scm-step:2.11
workflow-step-api:2.23
workflow-support:3.7`
)

func TestCanParseUpdateCentre_SingleDependency(t *testing.T) {
	u := update.Updater{}

	http := &api.FakeHTTP{}
	client := api.NewClient(api.ReplaceTripper(http))
	u.SetClient(client)
	u.IncludeDependencies()

	http.StubWithFixture(200, "golden.formatted.json")

	depsIn, err := update.FromStrings([]string{"scm-api"})
	assert.NoError(t, err)
	deps, err := u.LatestVersions(depsIn)
	assert.NoError(t, err)

	assert.Equal(t, len(http.Requests), 1)
	assert.Equal(t, http.Requests[0].URL.Path, "/update-center.actual.json")
	assert.Equal(t, http.Requests[0].URL.Host, "updates.jenkins.io")
	assert.Equal(t, http.Requests[0].URL.RawQuery, "")

	assert.Equal(t, 2, len(deps))
	assertContains(t, deps, "scm-api")
	assertContains(t, deps, "structs")
	assertContainsWithVersion(t, deps, "structs", "1.20")
}

func TestCanParseUpdateCentre_MultipleDependencies(t *testing.T) {
	u := update.Updater{}

	http := &api.FakeHTTP{}
	client := api.NewClient(api.ReplaceTripper(http))
	u.SetClient(client)
	u.IncludeDependencies()

	http.StubWithFixture(200, "golden.formatted.json")

	depsIn, err := update.FromStrings([]string{"scm-api", "job-dsl"})
	assert.NoError(t, err)
	deps, err := u.LatestVersions(depsIn)
	assert.NoError(t, err)

	assert.Equal(t, len(http.Requests), 1)
	assert.Equal(t, http.Requests[0].URL.Path, "/update-center.actual.json")
	assert.Equal(t, http.Requests[0].URL.Host, "updates.jenkins.io")
	assert.Equal(t, http.Requests[0].URL.RawQuery, "")

	assert.Equal(t, 4, len(deps))
	assertContains(t, deps, "scm-api")
	assertContains(t, deps, "structs")
	assertContainsWithVersion(t, deps, "structs", "1.20")
	assertContainsWithVersion(t, deps, "script-security", "1.54")
}

func TestCanParseUpdateCentre_CompleteList_NoVersions(t *testing.T) {
	u := update.Updater{}

	http := &api.FakeHTTP{}
	client := api.NewClient(api.ReplaceTripper(http))
	u.SetClient(client)

	http.StubWithFixture(200, "golden.full.json")

	depsIn, err := update.FromStrings([]string{
		"ansicolor",
		"antisamy-markup-formatter",
		"authentication-tokens",
		"azure-container-agents",
		"azure-vm-agents",
		"basic-branch-build-strategies",
		"blueocean",
		"blueocean-autofavorite",
		"blueocean-commons",
		"blueocean-config",
		"blueocean-core-js",
		"blueocean-dashboard",
		"blueocean-display-url",
		"blueocean-events",
		"blueocean-git-pipeline",
		"blueocean-github-pipeline",
		"blueocean-i18n",
		"blueocean-jira",
		"blueocean-jwt",
		"blueocean-personalization",
		"blueocean-pipeline-api-impl",
		"blueocean-pipeline-editor",
		"blueocean-pipeline-scm-api",
		"blueocean-rest",
		"blueocean-rest-impl",
		"blueocean-web",
		"branch-api",
		"build-name-setter",
		"config-file-provider",
		"cloud-stats",
		"configuration-as-code",
		"credentials",
		"credentials-binding",
		"datadog",
		"ec2",
		"embeddable-build-status",
		"extended-read-permission",
		"git",
		"git-client",
		"github",
		"github-api",
		"github-branch-source",
		"github-checks",
		"github-label-filter",
		"groovy",
		"inline-pipeline",
		"javadoc",
		"jira",
		"job-dsl",
		"junit",
		"kubernetes",
		"kubernetes-credentials",
		"kubernetes-credentials-provider",
		"ldap",
		"lockable-resources",
		"pipeline-utility-steps",
		"metrics",
		"matrix-auth",
		"matrix-project",
		"pipeline-build-step",
		"pipeline-github",
		"pipeline-graph-analysis",
		"pipeline-input-step",
		"pipeline-milestone-step",
		"pipeline-model-api",
		"pipeline-model-definition",
		"pipeline-model-extensions",
		"pipeline-rest-api",
		"pipeline-stage-step",
		"pipeline-stage-tags-metadata",
		"pipeline-stage-view",
		"plain-credentials",
		"prometheus",
		"scm-api",
		"scm-filter-branch-pr",
		"script-security",
		"ssh-agent",
		"ssh-credentials",
		"support-core",
		"timestamper",
		"token-macro",
		"toolenv",
		"variant",
		"warnings-ng",
		"workflow-aggregator",
		"workflow-api",
		"workflow-basic-steps",
		"workflow-cps",
		"workflow-cps-global-lib",
		"workflow-durable-task-step",
		"workflow-job",
		"workflow-multibranch",
		"workflow-scm-step",
		"workflow-step-api",
		"workflow-support"})
	assert.NoError(t, err)

	deps, err := u.LatestVersions(depsIn)
	assert.NoError(t, err)

	assert.Equal(t, 95, len(deps))
	assertContains(t, deps, "scm-api")
	assertContains(t, deps, "workflow-support")

	t.Log(update.AsStrings(deps))
}

func TestCanParseUpdateCentre_CompleteList_WithVersions(t *testing.T) {
	u := update.Updater{}

	http := &api.FakeHTTP{}
	client := api.NewClient(api.ReplaceTripper(http))
	u.SetClient(client)

	http.StubWithFixture(200, "golden.full.json")

	depsIn, err := update.FromStrings(strings.Split(versionSet, "\n"))
	assert.NoError(t, err)

	assert.NoError(t, err)
	deps, err := u.LatestVersions(depsIn)
	assert.NoError(t, err)

	assert.Equal(t, 95, len(deps))
	assertContains(t, deps, "scm-api")
	assertContains(t, deps, "workflow-support")

	assertContainsWithVersion(t, deps, "ansicolor", "0.7.3")

	t.Log(update.AsStrings(deps))
}

func TestCanParseUpdateCentre_CompleteList_CanDetectUpdates(t *testing.T) {
	u := update.Updater{}

	http := &api.FakeHTTP{}
	client := api.NewClient(api.ReplaceTripper(http))
	u.SetClient(client)

	http.StubWithFixture(200, "golden.full.json")

	updatedVersionSet := strings.ReplaceAll(versionSet, "ansicolor:0.7.3", "ansicolor:0.0.1")
	depsIn, err := update.FromStrings(strings.Split(updatedVersionSet, "\n"))
	assert.NoError(t, err)

	assert.NoError(t, err)
	deps, err := u.LatestVersions(depsIn)
	assert.NoError(t, err)

	assert.Equal(t, 95, len(deps))
	assertContains(t, deps, "scm-api")
	assertContains(t, deps, "workflow-support")

	assertContainsWithVersion(t, deps, "ansicolor", "0.7.3")

	t.Log(update.AsStrings(deps))

	changed := update.FindAll(deps, func(info update.DepInfo) bool {
		return info.Changed
	})

	t.Log(update.AsStrings(changed))
	assert.Equal(t, 1, len(changed))
}

func assertContains(t *testing.T, deps []update.DepInfo, name string) {
	for _, a := range deps {
		if a.Name == name {
			return
		}
	}
	assert.Fail(t, fmt.Sprintf("Was expecting to see %s", name))
}

func assertContainsWithVersion(t *testing.T, deps []update.DepInfo, name string, version string) {
	for _, a := range deps {
		if a.Name == name && a.Version == version {
			return
		}
	}
	assert.Fail(t, fmt.Sprintf("Was expecting to see %s:%s", name, version))
}
