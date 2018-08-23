package semver_test

import (
	"testing"

	"github.com/coreos/go-semver/semver"
	semverutil "github.com/platform9/ssh-provider/pkg/util/semver"
)

func TestEqualMajorMinorPatchVersions(t *testing.T) {
	tcs := []struct {
		name  string
		a     string
		b     string
		equal bool
	}{
		{
			name:  "equal",
			a:     "0.0.2-9+8d7d5693ad4ec9",
			b:     "0.0.2-10+g61d9a1a",
			equal: true,
		},
		{
			name:  "unequal",
			a:     "0.0.2-9+8d7d5693ad4ec9",
			b:     "0.0.3-9+8d7d5693ad4ec9",
			equal: false,
		},
	}

	for _, tc := range tcs {
		a := semver.New(tc.a)
		b := semver.New(tc.b)
		if semverutil.EqualMajorMinorPatchVersions(*a, *b) != tc.equal {
			t.Errorf("%s and %s should be equal, ignoring their pre-release identifiers", a, b)
		}
	}
}
