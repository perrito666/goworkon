package goinstalls

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// GODLURL is the url from which golang src tars can be dowloaded.
const GODLURL = "https://storage.googleapis.com/golang/"

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

func versionFromString(versionStr string) (Version, error) {
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

const srcPrefix = ".src.tar.gz"

// OnlineAvailableVersions returns a map of all found versions grouped
// by Minor number and with the latest patch of said Minor as value.
func OnlineAvailableVersions() (map[Version]string, error) {
	response, err := http.Get(GODLURL)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer response.Body.Close()
	var rawXML bytes.Buffer
	_, err = io.Copy(&rawXML, response.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	//fmt.Println(rawXML.String())
	var downloadables golangDownloadsPageContents
	if err := xml.Unmarshal(rawXML.Bytes(), &downloadables); err != nil {
		return nil, errors.WithStack(err)
	}
	versions := map[Version]string{}
	for _, dl := range downloadables.Contents {
		if strings.HasSuffix(dl.Key, "src.tar.gz") {
			versionString := dl.Key[2 : len(dl.Key)-len(srcPrefix)]
			v, err := versionFromString(versionString)
			if err != nil {
				continue
			}
			versions[v] = dl.Key
		}
	}
	// TODO (perrito666) Return only the newest of each minor.
	return versions, nil
}

// Update gets goVersion to the lates version of the
// given major.
func Update(goVersion string, updateRepos bool) {
}
