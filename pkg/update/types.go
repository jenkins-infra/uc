package update

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var semverRE = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)

type DepInfo struct {
	Name    string
	Version string
	Comment string
	Changed bool
}

func (d *DepInfo) String() string {
	v := d.versionOrEmpty()
	if v == "" {
		return fmt.Sprintf("%s%s", d.Name, d.formattedComment())
	}
	return fmt.Sprintf("%s:%s%s", d.Name, d.versionOrEmpty(), d.formattedComment())
}

func (d *DepInfo) versionOrEmpty() string {
	if d.Version == "0.0.0" {
		return ""
	}
	return d.Version
}

func (d *DepInfo) formattedComment() string {
	if d.Comment == "" {
		return ""
	}
	return fmt.Sprintf(" # %s", d.Comment)
}

func (d *DepInfo) ShouldUpdate(version string) bool {
	if isSemverCompatible(version) && isSemverCompatible(d.Version) {
		return compareSemvers(d.Version, version)
	} else {
		return compareNonSemvers(d.Version, version)
	}

	return false;
}

func compareSemvers(v1string string, v2string string) bool {
	v1 := semver.MustParse(v1string)
	v2 := semver.MustParse(v2string)

	if v1 != nil && v2 != nil {
		if v2.GreaterThan(v1) {
			return true
		}
	}
	
	return false
}

func compareNonSemvers(v1string string, v2string string) bool {
	// lets split these versions by the periods and compare each part
	parts1 := strings.Split(v1string, ".")
	parts2 := strings.Split(v2string, ".")

	for i := 0; i < max(len(parts1), len(parts2)); i++ {
		part1 := safePart(i, parts1)
		part2 := safePart(i, parts2)
		logrus.Debugf("comparing '%s' and '%s'", part1, part2)
		if part1 != part2 {
			return strings.Compare(part1, part2) < 0
		}
	}

	return false
}

func safePart(index int, parts []string) string {
	if index < len(parts)  {
		return parts[index]
	}
	return ""
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func isSemverCompatible(version string) bool {
	return semverRE.MatchString(version)
}

func (d *DepInfo) SkipUpdate() bool {
	return strings.Contains(d.Comment, "noupdate")
}

func FromString(in string) (*DepInfo, error) {
	di := DepInfo{}

	if strings.Contains(in, "#") {
		parts := strings.Split(in, "#")
		if len(parts) != 2 {
			return nil, errors.New("unable to parse comment for " + in)
		}
		di.Comment = strings.TrimSpace(parts[1])
		in = strings.TrimSpace(parts[0])
	}

	if strings.Contains(in, ":") {
		parts := strings.Split(in, ":")
		if len(parts) != 2 {
			return nil, errors.New("unable to parse plugin:version for " + in)
		}
		di.Name = parts[0]
		di.Version = parts[1]
		return &di, nil
	}

	di.Name = in
	di.Version = "0.0.0"
	return &di, nil
}

func FromStrings(input []string) ([]DepInfo, error) {
	deps := []DepInfo{}
	for _, in := range input {
		d, err := FromString(in)
		if err != nil {
			return nil, err
		}
		deps = append(deps, *d)
	}
	return deps, nil
}

func AsStrings(deps []DepInfo) []string {
	out := []string{}
	for _, d := range deps {
		out = append(out, d.String())
	}
	return out
}

func FindAll(deps []DepInfo, test func(info DepInfo) bool) (ret []DepInfo) {
	for _, d := range deps {
		if test(d) {
			ret = append(ret, d)
		}
	}
	return
}
