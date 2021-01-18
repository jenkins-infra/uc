package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/garethjevans/uc/pkg/update"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// UpdateCmd defines the cmd.
type UpdateCmd struct {
	common

	Write               bool
	IncludeDependencies bool
	DisplayUpdates      bool
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

	c.Cmd = cmd
	c.addCommonFlags()

	cmd.Flags().BoolVarP(&c.Write, "write", "w", false,
		"Update the file rather than display to stdout")
	cmd.Flags().BoolVarP(&c.IncludeDependencies, "include-dependencies", "d", false,
		"Add any additional dependencies to the output")
	cmd.Flags().BoolVarP(&c.DisplayUpdates, "display-updates", "u", false,
		"Write updates to stdout")

	return cmd
}

// Run update help.
func (c *UpdateCmd) Run() error {
	err := c.validateCommonFlags()
	if err != nil {
		return err
	}

	depsIn, err := c.readFromPath()
	if err != nil {
		return err
	}

	if c.Updater == nil {
		c.Updater = &update.Updater{}
		version, err := c.determineVersion()
		if err != nil {
			return err
		}
		c.Updater.SetVersion(version)
	}

	if c.IncludeDependencies {
		c.Updater.IncludeDependencies()
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
