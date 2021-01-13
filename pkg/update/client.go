package update

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/garethjevans/updatecenter/pkg/api"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type DepInfo struct {
	Name    string
	Version string
}

func (d *DepInfo) String() string {
	return fmt.Sprintf("%s:%s", d.Name, d.Version)
}

func FromString(in string) (*DepInfo, error) {
	if strings.Contains(in, ":") {
		parts := strings.Split(in, ":")
		if len(parts) != 2 {
			return nil, errors.New("unable to parse plugin:version for " + in)
		}
		return &DepInfo{Name: parts[0], Version: parts[1]}, nil
	}
	return &DepInfo{Name: in, Version: "0.0.0"}, nil
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

func AsStrings(deps []DepInfo) []string{
	out := []string{}
	for _, d := range deps {
		out = append(out, d.String())
	}
	return out
}


type Updater struct {
	config              *Config
	client              *api.Client
	version             string
	includeDependencies bool
}

func (u *Updater) IncludeDependencies() {
	u.includeDependencies = true
}

func (u *Updater) SetVersion(version string) {
	u.version = version
}

func (u *Updater) SetClient(client *api.Client) {
	u.client = client
}

func (u *Updater) Client() api.Client {
	if u.client == nil {
		u.client = api.BasicClient()
	}
	return *u.client
}

func (u *Updater) get() error {
	if u.config == nil {
		c := &Config{}

		err := u.Client().GET(u.version, c)
		if err != nil {
			return err
		}

		u.config = c
	}
	return nil
}

func (u *Updater) LatestVersions(plugins []DepInfo) ([]DepInfo, error) {
	if u.config == nil {
		err := u.get()
		if err != nil {
			return nil, err
		}
	}

	deps := []DepInfo{}
	for _, p := range u.config.Plugins {
		if Contains(plugins, p.Name) {
			// add the plugin
			if !Contains(deps, p.Name) {
				deps = append(deps, DepInfo{Name: p.Name, Version: p.Version})
			} else {
				err := setVersionIfNewer(deps, p.Name, p.Version)
				if err != nil {
					return nil, err
				}
			}

			if u.includeDependencies {
				// add the plugin dependencies
				for _, d := range p.Dependencies {
					if !d.Optional {
						if !Contains(deps, d.Name) {
							deps = append(deps, DepInfo{Name: d.Name, Version: d.Version})
						} else {
							err := setVersionIfNewer(deps, d.Name, d.Version)
							if err != nil {
								return nil, err
							}
						}
					}
				}
			}
		}
	}

	sort.Slice(deps, func(i, j int) bool {
		return deps[i].Name < deps[j].Name
	})

	return deps, nil
}

func setVersionIfNewer(deps []DepInfo, name string, version string) error {
	for i := range deps {
		if deps[i].Name == name {
			v1, err := semver.NewVersion(deps[i].Version)
			if err != nil {
				logrus.Debugf("version %s is invalid for %s", deps[i].Version, name)
			}

			v2, err := semver.NewVersion(version)
			if err != nil {
				logrus.Debugf("version %s is invalid for %s", version, name)
			}

			if v1 != nil && v2 != nil {
				if v2.GreaterThan(v1) {
					deps[i].Version = version
				}
			}
		}
	}
	return nil
}
