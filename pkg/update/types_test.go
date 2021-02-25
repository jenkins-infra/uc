package update_test

import (
	"testing"

	"github.com/garethjevans/uc/pkg/update"
	"github.com/stretchr/testify/assert"
)

func TestFromString(t *testing.T) {
	tests := []struct {
		in      string
		name    string
		version string
		comment string
	}{
		{in: "ldap", name: "ldap", version: "0.0.0"},
		{in: "ldap:1.26", name: "ldap", version: "1.26"},
		{in: "ldap:1.27", name: "ldap", version: "1.27"},
		{in: "ldap:1.0 # this is a comment", name: "ldap", version: "1.0", comment: "this is a comment"},
		{in: "ldap:1.0 # noupdate", name: "ldap", version: "1.0", comment: "noupdate"},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			di, err := update.FromString(tc.in)
			assert.NoError(t, err)
			assert.Equal(t, di.Name, tc.name)
			assert.Equal(t, di.Version, tc.version)

			assert.Equal(t, di.Comment, tc.comment)
			assert.Equal(t, di.String(), tc.in)
		})
	}
}

func TestFromStringSkipUpdate(t *testing.T) {
	tests := []struct {
		in         string
		skipUpdate bool
	}{
		{in: "ldap", skipUpdate: false},
		{in: "ldap:1.26", skipUpdate: false},
		{in: "ldap:1.27", skipUpdate: false},
		{in: "ldap:1.0 # this is a comment", skipUpdate: false},
		{in: "ldap:1.0 # noupdate", skipUpdate: true},
	}

	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			di, err := update.FromString(tc.in)
			assert.NoError(t, err)
			assert.Equal(t, di.SkipUpdate(), tc.skipUpdate)
			assert.Equal(t, di.String(), tc.in)
		})
	}
}
