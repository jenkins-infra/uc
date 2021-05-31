package cmd_test

import (
	"github.com/jenkins-infra/uc/pkg/cmd"
	"path/filepath"
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
	// antisamy-markup-formatter:2.1
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
	// antisamy-markup-formatter:2.1
	// authentication-tokens:1.4
}