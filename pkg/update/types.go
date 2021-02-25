package update

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

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
