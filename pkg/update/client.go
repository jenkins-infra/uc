package update

import (
	"github.com/Masterminds/semver"
	"github.com/garethjevans/updatecenter/pkg/api"
)

type DepInfo struct {
	Name    string
	Version string
}

type Updater struct {
	config *Config
	client *api.Client
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

		err := u.Client().GET("", c)
		if err != nil {
			return err
		}

		u.config = c
	}
	return nil
}

func (u *Updater) LatestVersions(plugins []string) ([]DepInfo, error) {
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
			if !contains(deps, p.Name) {
				deps = append(deps, DepInfo{Name: p.Name, Version: p.Version})
			} else {
				err := setVersionIfNewer(deps, p.Name, p.Version)
				if err != nil {
					return nil, err
				}
			}

			// add the plugin dependencies
			for _, d := range p.Dependencies {
				if !d.Optional {
					if !contains(deps, d.Name) {
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

	return deps, nil
}

func contains(deps []DepInfo, name string) bool {
	for _, d := range deps {
		if d.Name == name {
			return true
		}
	}
	return false
}

func setVersionIfNewer(deps []DepInfo, name string, version string) error {
	for i := range deps {
		if deps[i].Name == name {
			v1, err := semver.NewVersion(deps[i].Version)
			if err != nil {
				return err
			}

			v2, err := semver.NewVersion(version)
			if err != nil {
				return err
			}

			if v2.GreaterThan(v1) {
				deps[i].Version = version
			}
		}
	}
	return nil
}
