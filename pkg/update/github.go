package update

import (
	"regexp"

	"github.com/jenkins-infra/uc/pkg/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type GitHub struct {
	client *api.Client
}

func (g *GitHub) SetClient(client *api.Client) {
	g.client = client
}

func (g *GitHub) Client() api.Client {
	if g.client == nil {
		g.client = api.BasicClient()
	}
	return *g.client
}

func (g *GitHub) GetLatestLTSRelease() (string, error) {
	releaseInfo := []ReleaseInfo{}
	err := g.Client().REST("https://api.github.com/repos/jenkinsci/jenkins/releases?per_page=100", nil, &releaseInfo)
	if err != nil {
		return "", err
	}

	logrus.Debugf("got %d records", len(releaseInfo))

	ltsReleaseInfo := filterWithRegexp(releaseInfo, `^([\d]+)\.([\d]+)(\.[\d]+)`)

	logrus.Debugf("filtered down to %d lts records", len(ltsReleaseInfo))

	if len(ltsReleaseInfo) == 0 {
		return "", errors.New("unable to determine latest lts release from github")
	}

	return ltsReleaseInfo[0].Name, nil
}

type ReleaseInfo struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	Name            string `json:"name"`
	Draft           bool   `json:"draft"`
	Prerelease      bool   `json:"prerelease"`
	CreatedAt       string `json:"created_at"`
	PublishedAt     string `json:"published_at"`
}

func filterWithRegexp(in []ReleaseInfo, match string) (ret []ReleaseInfo) {
	for _, i := range in {
		r, _ := regexp.Compile(match)
		if r.MatchString(i.Name) {
			ret = append(ret, i)
		}
	}
	return
}
