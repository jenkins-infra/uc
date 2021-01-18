package cmd

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/garethjevans/uc/pkg/update"
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
	logrus.Debugf("reading plugins from %s", c.Path)
	data, err := ioutil.ReadFile(c.Path)
	if err != nil {
		return nil, errors.New("unable to read file from path " + c.Path)
	}

	lines := strings.Split(string(data), "\n")
	depsIn, err := update.FromStrings(lines)
	if err != nil {
		return nil, errors.New("unable to convert file into dependencies")
	}

	return depsIn, nil
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
