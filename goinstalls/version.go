package goinstalls

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Version represents a go version, since its intended for
// the source of released versions, no suffix/build is supported.
// TODO (perrito666) accept build/rc and figure how to determine
// newer.
type Version struct {
	// Major is the first tier of a version number.
	Major int
	// Minor is the second tier of a version number.
	Minor int
	// Patch is the final tier of a version number.
	Patch int
}

// IsNewerThan returns true if the passed version is
// older than the current one.
func (v Version) IsNewerThan(ver Version) bool {
	if v.Major > ver.Major {
		return true
	}
	if v.Major < ver.Major {
		return false
	}

	if v.Minor > ver.Minor {
		return true
	}
	if v.Minor < ver.Minor {
		return false
	}

	if v.Patch > ver.Patch {
		return true
	}
	return false
}

// SameVersion returns true if the passed version is
// of the same root than the current one, ie: 1.7.2 and 1.7.3.
func (v Version) SameVersion(ver Version) bool {
	return v.Major == ver.Major && v.Minor == ver.Minor
}

// CommonVersionString returns a string representing the
// Major.Minor version of the version.
func (v Version) CommonVersionString() string {
	return fmt.Sprintf("%d.%d", v.Major, v.Minor)
}

// String returns the string representation of this Version.
func (v Version) String() string {
	if v.Patch == 0 {
		return v.CommonVersionString()
	}
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

type downloadable struct {
	Key            string `xml:"Key"`
	Generation     string `xml:"Generation"`
	MetaGeneration int64  `xml:"MetaGeneration"`
	LastModified   string `xml:"LastModified"`
	ETag           string `xml:"ETag"`
	Size           int64  `xml:"Size"`
	Owner          string `xml:"Owner"`
}

type golangDownloadsPageContents struct {
	Name     string         `xml:"Name"`
	Contents []downloadable `xml:"Contents"`
}

type golangDownloadsPage struct {
	ListBucketResult golangDownloadsPageContents `xml:"ListBucketResult"`
}

// VersionFromString returns a Version created from a valid string version
// of the form x.y.z or x.y
func VersionFromString(versionStr string) (Version, error) {
	versionParts := strings.Split(versionStr, ".")

	if len(versionParts) < 2 || len(versionParts) > 3 {
		return Version{}, errors.Errorf("%q is not a valid version string", versionStr)
	}

	major, err := strconv.Atoi(versionParts[0])
	if err != nil {
		return Version{}, errors.Wrapf(err, "%q is an invalid major", versionParts[0])
	}
	minor, err := strconv.Atoi(versionParts[1])
	if err != nil {
		// This is an rc/beta most likely.
		return Version{}, errors.Wrapf(err, "%q is an invalid minor", versionParts[1])
	}

	if len(versionParts) == 2 {
		return Version{
			Major: major,
			Minor: minor,
			Patch: 0,
		}, nil
	}

	if len(versionParts) == 3 {
		patch, err := strconv.Atoi(versionParts[2])
		if err != nil {
			return Version{}, errors.Wrapf(err, "%q is an invalid patch", versionParts[2])
		}
		return Version{
			Major: major,
			Minor: minor,
			Patch: patch,
		}, nil
	}
	return Version{}, errors.Errorf("could not parse version string %q", versionStr)
}
