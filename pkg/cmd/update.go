package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/garethjevans/uc/pkg/update"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// UpdateCmd defines the cmd.
type UpdateCmd struct {
	Cmd  *cobra.Command
	Args []string

	Updater                        *update.Updater
	Path                           string
	JenkinsVersion                 string
	Write                          bool
	IncludeDependencies            bool
	DisplayUpdates                 bool
	DetermineVersionFromDockerfile bool
	DockerFilePath                 string
}

// NewUpdateCmd defines a new cmd.
func NewUpdateCmd() *cobra.Command {
	c := &UpdateCmd{}
	cmd := &cobra.Command{
		Use:   "update",
		Short: "uc update --path <path>",
		Long: `To update all plugins against the latest version of Jenkins:

    uc update --path <path>

To update all plugins against a specific version of Jenkins:

    uc update --path <path> --jenkins-version <version>
`,
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {
			c.Cmd = cmd
			c.Args = args
			err := c.Run()
			if err != nil {
				logrus.Errorf("unhandled error - %s", err)
				logrus.Fatal("unable to run command")
			}
		},
	}

	cmd.Flags().StringVarP(&c.Path, "path", "p", "plugins.txt",
		"Path to the plugins.txt file")
	cmd.Flags().StringVarP(&c.JenkinsVersion, "jenkins-version", "j", "",
		"The version of Jenkins to query against")
	cmd.Flags().BoolVarP(&c.Write, "write", "w", false,
		"Update the file rather than display to stdout")
	cmd.Flags().BoolVarP(&c.IncludeDependencies, "include-dependencies", "d", false,
		"Add any additional dependencies to the output")
	cmd.Flags().BoolVarP(&c.DisplayUpdates, "display-updates", "u", false,
		"Write updates to stdout")
	cmd.Flags().BoolVarP(&c.DetermineVersionFromDockerfile, "determine-version-from-dockerfile", "", false,
		"Attempt to determine the Jenkins version from a Dockerfile")
	cmd.Flags().StringVarP(&c.DockerFilePath, "dockerfile-path", "", "Dockerfile",
		"Path to the Dockerfile")

	return cmd
}

// Run update help.
func (c *UpdateCmd) Run() error {
	if c.Path == "" {
		return errors.New("--path needs to be set")
	}

	if c.DetermineVersionFromDockerfile && c.JenkinsVersion != "" {
		return errors.New("only one of --determine-version-from-dockerfile or --jenkins-version should be used")
	}

	data, err := ioutil.ReadFile(c.Path)
	if err != nil {
		return errors.New("unable to read file from path " + c.Path)
	}

	lines := strings.Split(string(data), "\n")
	depsIn, err := update.FromStrings(lines)
	if err != nil {
		return errors.New("unable to convert file into dependencies")
	}

	if c.Updater == nil {
		c.Updater = &update.Updater{}
	}

	if c.IncludeDependencies {
		c.Updater.IncludeDependencies()
	}

	if c.DetermineVersionFromDockerfile {
		r, err := os.Open(c.DockerFilePath)
		if err != nil {
			return errors.New("unable to open " + c.DockerFilePath)
		}
		fullVersion, err := update.DetermineJenkinsVersionFromDockerfile(r)
		if err != nil {
			return errors.New("unable to determine version from Dockerfile")
		}
		c.Updater.SetVersion(update.ExtractExactVersion(fullVersion))
	} else if c.JenkinsVersion != "" {
		c.Updater.SetVersion(c.JenkinsVersion)
	}

	depsOut, err := c.Updater.LatestVersions(depsIn)
	if err != nil {
		return errors.New("unable to determine latest versions")
	}

	// filter out any empty records
	depsOut = update.FindAll(depsOut, func(info update.DepInfo) bool {
		return info.Name != ""
	})

	depsString := update.AsStrings(depsOut)
	dataToWrite := strings.Join(depsString, "\n")

	changed := update.FindAll(depsOut, func(info update.DepInfo) bool {
		return info.Changed
	})

	changedString := update.AsStrings(changed)

	if c.Write {
		bytesToWrite := []byte(dataToWrite)
		err := ioutil.WriteFile(c.Path, bytesToWrite, 0600)
		if err != nil {
			return errors.New("unable to write file")
		}
		if c.DisplayUpdates {
			fmt.Println(strings.Join(changedString, "\n"))
		}
	} else {
		if c.DisplayUpdates {
			fmt.Println(strings.Join(changedString, "\n"))
		} else {
			fmt.Println(dataToWrite)
		}
	}

	return nil
}
