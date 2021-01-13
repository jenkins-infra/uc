package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/garethjevans/updatecenter/pkg/update"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// UpdateCmd defines the cmd.
type UpdateCmd struct {
	Cmd  *cobra.Command
	Args []string

	Updater             *update.Updater
	Path                string
	JenkinsVersion      string
	Write               bool
	IncludeDependencies bool
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
				logrus.Fatal("unable to run command")
			}
		},
	}

	cmd.Flags().StringVarP(&c.Path, "path", "p", "",
		"Path to the plugins.txt file")
	cmd.Flags().StringVarP(&c.JenkinsVersion, "jenkins-version", "j", "",
		"The version of Jenkins to query against")
	cmd.Flags().BoolVarP(&c.Write, "write", "w", false,
		"Update the file rather than display to stdout")
	cmd.Flags().BoolVarP(&c.IncludeDependencies, "include-dependencies", "d", false,
		"Add any additional dependencies to the output")

	return cmd
}

// Run update help.
func (c *UpdateCmd) Run() error {
	if c.Path == "" {
		return errors.New("--path needs to be set")
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

	depsOut, err := c.Updater.LatestVersions(depsIn)
	if err != nil {
		return errors.New("unable to determine latest versions")
	}

	depsString := []string{}
	for _, dep := range depsOut {
		depsString = append(depsString, dep.String())
	}

	dataToWrite := strings.Join(depsString, "\n")

	if c.Write {
		bytesToWrite := []byte(dataToWrite)
		err := ioutil.WriteFile(c.Path, bytesToWrite, 0600)
		if err != nil {
			return errors.New("unable to write file")
		}
	} else {
		fmt.Println(dataToWrite)
	}

	return nil
}
