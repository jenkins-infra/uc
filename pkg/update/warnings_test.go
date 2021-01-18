package update_test

import (
	"testing"

	"github.com/garethjevans/uc/pkg/update"
	"github.com/stretchr/testify/assert"
)

func TestCanEvaluateWarnings(t *testing.T) {
	w := update.WarningInfo{
		Name: "google-login",
		Versions: []update.VersionInfo{
			{
				LastVersion: "1.1",
				Pattern:     `1[.][01](|[.-].*)`,
			},
		},
	}

	assert.True(t, w.Matches("1.0"))
	assert.True(t, w.Matches("1.1"))
	assert.False(t, w.Matches("1.2"))
}
