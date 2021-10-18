package cmd_test

import (
	"path/filepath"

	"github.com/jenkins-infra/uc/pkg/cmd"
)

func ExampleDisplayUpdatesFromPluginTxt() {
	c := cmd.UpdateCmd{}
	c.Path = filepath.Join("testdata", "plugins.txt")
	err := c.Run()
	if err != nil {
		panic(err)
	}

	// Output:
	// ansicolor:1.0.0
	// antisamy-markup-formatter:2.3
	// authentication-tokens:1.4
}

func ExampleDisplayUpdatesFromValuesYaml() {
	c := cmd.UpdateCmd{}
	c.Path = filepath.Join("testdata", "values.yaml")
	c.YamlLocation = "controller.installPlugins"
	err := c.Run()
	if err != nil {
		panic(err)
	}

	// Output:
	// ansicolor:1.0.0
	// antisamy-markup-formatter:2.3
	// authentication-tokens:1.4
}
