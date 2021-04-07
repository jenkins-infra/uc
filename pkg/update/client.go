package update

import (
	"sort"

	"github.com/garethjevans/uc/pkg/api"
	"github.com/sirupsen/logrus"
)

type Updater struct {
	config              *Config
	client              *api.Client
	version             string
	includeDependencies bool
	securityUpdates     bool
}

func (u *Updater) IncludeDependencies() {
	u.includeDependencies = true
}

func (u *Updater) SecurityUpdates() {
	u.securityUpdates = true
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

	warnings, err := u.GetWarnings(plugins)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("got %d warning(s)", len(warnings))

	for _, plugin := range plugins {
		if !plugin.SkipUpdate() {
			p := u.config.Plugins[plugin.Name]
			if u.securityUpdates {
				logrus.Debugf("checking if there is a security update for %s", p.Name)
				if u.isSecurityUpdateForPlugin(warnings, p.Name) {
					logrus.Debugf("security update available for %s, update to %s", p.Name, p.Version)
					setVersionIfNewer(plugins, p.Name, p.Version)
				}
			} else {
				setVersionIfNewer(plugins, p.Name, p.Version)
			}

			if u.includeDependencies {
				plugins = addDependenciesForPlugin(plugins, p.Dependencies)
			}
		}
	}

	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].Name < plugins[j].Name
	})

	return plugins, nil
}

func (u *Updater) GetWarnings(plugins []DepInfo) ([]WarningInfo, error) {
	if u.config == nil {
		err := u.get()
		if err != nil {
			return nil, err
		}
	}

	warnings := []WarningInfo{}

	for _, p := range plugins {
		for _, w := range u.config.Warnings {
			if w.Name == p.Name {
				if w.Matches(p.Version) {
					logrus.Debugf("matches warning for %s & %s", w.Name, p.Version)
					warnings = append(warnings, w)
				}
			}
		}
	}

	return warnings, nil
}

func setVersionIfNewer(deps []DepInfo, name string, version string) {
	for i := range deps {
		if deps[i].Name == name {
			if deps[i].ShouldUpdate(version) {
				deps[i].Version = version
				deps[i].Changed = true
			}
		}
	}
}

func (u *Updater) isSecurityUpdateForPlugin(warnings []WarningInfo, plugin string) bool {
	for _, w := range warnings {
		if w.Name == plugin {
			return true
		}
	}
	return false
}

func addDependenciesForPlugin(deps []DepInfo, dependencies []Dependency) []DepInfo {
	// add the plugin dependencies
	for _, d := range dependencies {
		if !d.Optional {
			if !Contains(deps, d.Name) {
				deps = append(deps, DepInfo{Name: d.Name, Version: d.Version, Changed: true})
			} else {
				setVersionIfNewer(deps, d.Name, d.Version)
			}
		}
	}
	return deps
}
