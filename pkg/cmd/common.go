package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-yaml/yaml"

	"github.com/jenkins-infra/uc/pkg/update"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type common struct {
	Cmd  *cobra.Command
	Args []string

	Updater                        *update.Updater
	Path                           string
	JenkinsVersion                 string
	DetermineVersionFromDockerfile bool
	DockerFilePath                 string
	YamlLocation                   string
}

func (c *common) addCommonFlags() {
	c.Cmd.Flags().StringVarP(&c.Path, "path", "p", "plugins.txt",
		"Path to the plugins.txt file")
	c.Cmd.Flags().StringVarP(&c.JenkinsVersion, "jenkins-version", "j", "",
		"The version of Jenkins to query against")
	c.Cmd.Flags().BoolVarP(&c.DetermineVersionFromDockerfile, "determine-version-from-dockerfile", "", false,
		"Attempt to determine the Jenkins version from a Dockerfile")
	c.Cmd.Flags().StringVarP(&c.DockerFilePath, "dockerfile-path", "", "Dockerfile",
		"Path to the Dockerfile")
	c.Cmd.Flags().StringVarP(&c.YamlLocation, "yaml-location", "", "controller.installPlugins",
		"Location in the yaml file of the plugin declaration")
}

func (c *common) validateCommonFlags() error {
	if c.Path == "" {
		return errors.New("--path needs to be set")
	}

	if c.DetermineVersionFromDockerfile && c.JenkinsVersion != "" {
		return errors.New("only one of --determine-version-from-dockerfile or --jenkins-version should be used")
	}

	return nil
}

func (c *common) readFromPath() ([]update.DepInfo, error) {
	var depsIn []update.DepInfo
	if strings.HasSuffix(c.Path, ".txt") {
		logrus.Debugf("reading plugins from %s", c.Path)
		data, err := ioutil.ReadFile(c.Path)
		if err != nil {
			return nil, errors.New("unable to read file from path " + c.Path)
		}

		lines := strings.Split(string(data), "\n")
		depsIn, err = update.FromStrings(lines)
		if err != nil {
			return nil, errors.New("unable to convert file into dependencies")
		}
	} else if c.isYaml() {
		logrus.Debugf("Attempting to parse as YAML")
		data, err := ioutil.ReadFile(c.Path)
		if err != nil {
			return nil, errors.New("unable to read file from path " + c.Path)
		}

		var values map[interface{}]interface{}
		err = yaml.Unmarshal(data, &values)
		if err != nil {
			return nil, errors.New("unable to unmarshall")
		}

		logrus.Debugf("Values=%s", values)
		index := strings.Split(c.YamlLocation, ".")

		pluginsFromYaml, err := locatePluginListFromYaml(values, index)
		if err != nil {
			return nil, errors.New("unable to locate plugins from yaml")
		}
		depsIn, err = update.FromStrings(pluginsFromYaml)
		if err != nil {
			return nil, errors.New("unable to convert file into dependencies")
		}

		return depsIn, nil
	} else {
		return nil, errors.New("unsupported file type " + c.Path)
	}
	return depsIn, nil
}

func locatePluginListFromYaml(values map[interface{}]interface{}, index []string) ([]string, error) {
	for _, i := range index {
		logrus.Debugf("looking for %s", i)
		newValues, ok := values[i].(map[interface{}]interface{})
		if ok {
			logrus.Debugf("newValues=%s", newValues)
			values = newValues
		} else {
			pluginList, ok := values[i].([]interface{})
			if ok {
				stringPluginList := []string{}
				for _, plugin := range pluginList {
					stringPluginList = append(stringPluginList, fmt.Sprintf("%s", plugin))
				}
				return stringPluginList, nil
			}
			return nil, errors.New("unable to locate " + i)
		}
	}
	return nil, errors.New("unable to locate " + strings.Join(index, "."))
}

func (c *common) determineVersion() (string, error) {
	if c.DetermineVersionFromDockerfile {
		r, err := os.Open(c.DockerFilePath)
		if err != nil {
			return "", errors.New("unable to open " + c.DockerFilePath)
		}
		fullVersion, err := update.DetermineJenkinsVersionFromDockerfile(r)
		if err != nil {
			return "", errors.New("unable to determine version from Dockerfile")
		}
		extractedVersion := update.ExtractExactVersion(fullVersion)
		logrus.Debugf("using jenkins version %s from dockerfile %s", extractedVersion, c.DockerFilePath)
		return extractedVersion, nil
	} else if c.JenkinsVersion != "" {
		logrus.Debugf("using jenkins version %s", c.JenkinsVersion)
		return c.JenkinsVersion, nil
	}
	return "", nil
}

func (c *common) isYaml() bool {
	return strings.HasSuffix(c.Path, ".yaml") || strings.HasSuffix(c.Path, ".yml")
}
