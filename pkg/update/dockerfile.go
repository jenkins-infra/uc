package update

import (
	"io"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

func DetermineJenkinsVersionFromDockerfile(reader io.Reader) (string, error) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, reader)
	if err != nil {
		return "", err
	}

	lines := strings.Split(buf.String(), "\n")
	lines = filter(lines, func(l string) bool {
		return strings.HasPrefix(l, "FROM ")
	})

	if len(lines) != 1 {
		return "", errors.New("unable to determine parent image name")
	}

	from := lines[0]
	fullImageName := strings.ReplaceAll(from, "FROM ", "")

	if !strings.Contains(fullImageName, "jenkins/jenkins") {
		return "", errors.New("parent image does not appear to be a jenkins/jenkins image '" + fullImageName + "'")
	}

	if !strings.Contains(fullImageName, ":") {
		return "", errors.New("parent image does not appear to be tagged with a version '" + fullImageName + "'")
	}

	return strings.Split(fullImageName, ":")[1], nil
}

func IsLTS(version string) bool {
	return strings.Contains(version, "lts")
}

func ExtractExactVersion(jenkinsVersion string) string {
	r, _ := regexp.Compile(`^([\d]+)\.([\d]+)(?:\.[\d]+)?`)
	return r.FindString(jenkinsVersion)
}

func filter(in []string, test func(l string) bool) (ret []string) {
	for _, i := range in {
		if test(i) {
			ret = append(ret, i)
		}
	}
	return
}
