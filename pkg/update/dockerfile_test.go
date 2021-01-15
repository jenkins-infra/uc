package update_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/garethjevans/uc/pkg/update"
	"github.com/stretchr/testify/assert"
)

func TestDetermineJenkinsVersionFromDockerfile(t *testing.T) {
	testCases := []struct {
		fullVersionString string
		isLTS             bool
		expectedVersion   string
	}{
		{"jenkins/jenkins:jdk11-hotspot-windowsservercore-2019", false, ""},
		{"jenkins/jenkins:jdk11-hotspot-windowsservercore-1809", false, ""},
		{"jenkins/jenkins:2.263.2-lts-centos7", true, "2.263.2"},
		{"jenkins/jenkins:lts-centos7", true, ""},
		{"jenkins/jenkins:centos7", false, ""},
		{"jenkins/jenkins:2.263.2-lts-centos", true, "2.263.2"},
		{"jenkins/jenkins:lts-centos", true, ""},
		{"jenkins/jenkins:centos", false, ""},
		{"jenkins/jenkins:2.263.2-lts-jdk11", true, "2.263.2"},
		{"jenkins/jenkins:lts-jdk11", true, ""},
		{"jenkins/jenkins:jdk11", false, ""},
		{"jenkins/jenkins:2.263.2-lts-slim", true, "2.263.2"},
		{"jenkins/jenkins:lts-slim", true, ""},
		{"jenkins/jenkins:slim", false, ""},
		{"jenkins/jenkins:2.263.2-lts-alpine", true, "2.263.2"},
		{"jenkins/jenkins:lts-alpine", true, ""},
		{"jenkins/jenkins:alpine", false, ""},
		{"jenkins/jenkins:2.263.2-lts", true, "2.263.2"},
		{"jenkins/jenkins:lts", true, ""},
		{"jenkins/jenkins:latest", false, ""},
		{"jenkins/jenkins:2.275-centos7", false, "2.275"},
		{"jenkins/jenkins:2.275-centos", false, "2.275"},
		{"jenkins/jenkins:2.275-jdk11", false, "2.275"},
	}
	for _, tc := range testCases {
		t.Run(tc.fullVersionString, func(t *testing.T) {
			// #nosec G201
			reader := strings.NewReader(fmt.Sprintf("FROM %s\n\nENV E V\n", tc.fullVersionString))

			jenkinsVersion, err := update.DetermineJenkinsVersionFromDockerfile(reader)
			assert.NoError(t, err)

			assert.Equal(t, tc.isLTS, update.IsLTS(jenkinsVersion))
			assert.Equal(t, tc.expectedVersion, update.ExtractExactVersion(jenkinsVersion))
		})
	}
}
