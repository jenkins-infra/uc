package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/garethjevans/uc/pkg/update"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CheckCmd defines the cmd.
type CheckCmd struct {
	Cmd  *cobra.Command
	Args []string

	Updater                        *update.Updater
	Path                           string
	JenkinsVersion                 string
	DetermineVersionFromDockerfile bool
	DockerFilePath                 string
}

// NewCheckCmd defines a new cmd.
func NewCheckCmd() *cobra.Command {
	c := &CheckCmd{}
	cmd := &cobra.Command{
		Use:   "check",
		Short: "uc check --path <path>",
		Long: `Validate existing plugin versions against known vulnerabilities:

    uc check --path <path>

To update all plugins against a specific version of Jenkins:

    uc check --path <path> --jenkins-version <version>
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
	cmd.Flags().BoolVarP(&c.DetermineVersionFromDockerfile, "determine-version-from-dockerfile", "", false,
		"Attempt to determine the Jenkins version from a Dockerfile")
	cmd.Flags().StringVarP(&c.DockerFilePath, "dockerfile-path", "", "Dockerfile",
		"Path to the Dockerfile")

	return cmd
}

// Run update help.
func (c *CheckCmd) Run() error {
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

	warnings, err := c.Updater.GetWarnings(depsIn)
	if err != nil {
		return errors.New("unable to determine latest versions")
	}

	if len(warnings) == 0 {
		fmt.Println("No warning(s) found")
	} else {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Plugin", "Issue", "URL"})

		t.AppendSeparator()

		for _, w := range warnings {
			t.AppendRows([]table.Row{{w.Name, w.ID, w.URL}})
		}

		t.Render()
	}

	return nil
}
