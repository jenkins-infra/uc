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
		expectedVersion   string
		isLTS             bool
		expectedError     bool
	}{
		{fullVersionString: "jenkins/jenkins:jdk11-hotspot-windowsservercore-2019"},
		{fullVersionString: "jenkins/jenkins:jdk11-hotspot-windowsservercore-1809"},
		{fullVersionString: "jenkins/jenkins:2.263.2-lts-centos7", isLTS: true, expectedVersion: "2.263.2"},
		{fullVersionString: "jenkins/jenkins:lts-centos7", isLTS: true},
		{fullVersionString: "jenkins/jenkins:centos7"},
		{fullVersionString: "jenkins/jenkins:2.263.2-lts-centos", isLTS: true, expectedVersion: "2.263.2"},
		{fullVersionString: "jenkins/jenkins:centos"},
		{fullVersionString: "jenkins/jenkins:2.263.2-lts-jdk11", isLTS: true, expectedVersion: "2.263.2"},
		{fullVersionString: "jenkins/jenkins:lts-jdk11", isLTS: true},
		{fullVersionString: "jenkins/jenkins:jdk11"},
		{fullVersionString: "jenkins/jenkins:2.263.2-lts-slim", isLTS: true, expectedVersion: "2.263.2"},
		{fullVersionString: "jenkins/jenkins:lts-slim", isLTS: true},
		{fullVersionString: "jenkins/jenkins:slim"},
		{fullVersionString: "jenkins/jenkins:2.263.2-lts-alpine", isLTS: true, expectedVersion: "2.263.2"},
		{fullVersionString: "jenkins/jenkins:lts-alpine", isLTS: true},
		{fullVersionString: "jenkins/jenkins:alpine"},
		{fullVersionString: "jenkins/jenkins:2.263.2-lts", isLTS: true, expectedVersion: "2.263.2"},
		{fullVersionString: "jenkins/jenkins:lts", isLTS: true},
		{fullVersionString: "jenkins/jenkins:latest"},
		{fullVersionString: "jenkins/jenkins:2.275-centos7", expectedVersion: "2.275"},
		{fullVersionString: "jenkins/jenkins:2.275-centos", expectedVersion: "2.275"},
		{fullVersionString: "jenkins/jenkins:2.275-jdk11", expectedVersion: "2.275"},
		{fullVersionString: "random/baseimage:2.275-jdk11", expectedError: true},
		{fullVersionString: "random/baseimage", expectedError: true},
	}
	for _, tc := range testCases {
		t.Run(tc.fullVersionString, func(t *testing.T) {
			// #nosec G201
			reader := strings.NewReader(fmt.Sprintf("FROM %s\n\nENV E V\n", tc.fullVersionString))

			jenkinsVersion, err := update.DetermineJenkinsVersionFromDockerfile(reader)
			if tc.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.isLTS, update.IsLTS(jenkinsVersion))
				assert.Equal(t, tc.expectedVersion, update.ExtractExactVersion(jenkinsVersion))
			}
		})
	}
}
