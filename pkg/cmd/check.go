package cmd

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/garethjevans/uc/pkg/update"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CheckCmd defines the cmd.
type CheckCmd struct {
	common
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

	c.Cmd = cmd
	c.addCommonFlags()

	return cmd
}

// Run update help.
func (c *CheckCmd) Run() error {
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
