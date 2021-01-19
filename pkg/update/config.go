package update

import (
	"github.com/Masterminds/semver/v3"
	"github.com/Upliner/goback/regexp"
	"github.com/sirupsen/logrus"
)

// Config struct to store the update center config representation.
type Config struct {
	ConnectionCheckURL  string                     `json:"connectionCheckUrl"`
	Core                CoreInfo                   `json:"core"`
	Deprecations        map[string]DeprecationInfo `json:"deprecations"`
	ID                  string                     `json:"id"`
	Plugins             map[string]PluginInfo      `json:"plugins"`
	UpdateCenterVersion string                     `json:"updateCenterVersion"`
	Warnings            []WarningInfo              `json:"warnings"`
}

// CoreInfo struct to store the update center core representation.
type CoreInfo struct {
	BuildDate string `json:"buildDate"`
	Name      string `json:"core"`
	Sha1      string `json:"sha1"`
	Sha256    string `json:"sha256"`
	URL       string `json:"url"`
	Version   string `json:"version"`
}

// DeprecationInfo struct to store the update center deprecation representation.
type DeprecationInfo struct {
	URL string `json:"url"`
}

// PluginInfo struct to store the update center plugin representation.
type PluginInfo struct {
	BuildDate    string       `json:"buildDate"`
	Name         string       `json:"name"`
	Sha1         string       `json:"sha1"`
	Sha256       string       `json:"sha256"`
	URL          string       `json:"url"`
	Version      string       `json:"version"`
	RequiredCore string       `json:"requiredCore"`
	Dependencies []Dependency `json:"dependencies"`
}

// WarningInfo struct to store the update center warning representation.
type WarningInfo struct {
	ID       string        `json:"id"`
	Message  string        `json:"message"`
	Name     string        `json:"name"`
	Type     string        `json:"type"`
	URL      string        `json:"url"`
	Versions []VersionInfo `json:"versions"`
}

func (w *WarningInfo) Matches(in string) bool {
	for _, v := range w.Versions {
		if v.Matches(in) {
			return true
		}
	}
	return false
}

// VersionInfo struct to store the update center version representation.
type VersionInfo struct {
	LastVersion string `json:"lastVersion"`
	Pattern     string `json:"pattern"`
}

func (v *VersionInfo) Matches(in string) bool {
	r := regexp.MustCompile(v.Pattern)
	matches := r.MatchString(in)
	if !matches {
		return false
	}

	logrus.Debugf("matches - %s against %s", in, v.Pattern)
	logrus.Debugf("-> last version %s", v.LastVersion)

	lastVersion, err := semver.NewVersion(v.LastVersion)
	if err != nil {
		logrus.Debugf("lastVersion %s is invalid", lastVersion)
	}

	inVersion, err := semver.NewVersion(in)
	if err != nil {
		logrus.Infof("inVersion %s is invalid", in)
	}

	if lastVersion != nil && inVersion != nil {
		logrus.Debugf("checking   lastVersion %s >= inVersion %s", lastVersion, inVersion)
		if lastVersion.GreaterThan(inVersion) || lastVersion.Equal(inVersion) {
			return true
		}
	}

	return false
}

// Dependency struct to store the update center dependency representation.
type Dependency struct {
	Name     string `json:"name"`
	Optional bool   `json:"optional"`
	Version  string `json:"version"`
}
