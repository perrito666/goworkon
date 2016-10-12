package goinstalls

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// GODLURL is the url from which golang src tars can be dowloaded.
const GODLURL = "https://storage.googleapis.com/golang/"
const srcPrefix = ".src.tar.gz"

func filterNewer(vers map[Version]string) map[Version]string {
	unique := map[string]Version{}
	for k := range vers {
		s := k.CommonVersionString()
		v, ok := unique[s]
		if !ok {
			unique[s] = k
			continue
		}
		if k.IsNewerThan(v) {
			unique[s] = k
		}
	}
	result := make(map[Version]string, len(unique))
	for _, v := range unique {
		result[v] = vers[v]
	}
	return result
}

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
	return filterNewer(versions), nil
}

// InstalledAvailableVersions returns a slice of the Versions that
// have a current install locally.
func InstalledAvailableVersions() []Version {
	return []Version{}
}

func untar(tarFile *tar.Reader, targetPath string) error {
	for {
		h, err := tarFile.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.Wrap(err, "uncompressing headers")
		}
		p := filepath.Join(targetPath, h.Name)
		fmt.Println(p)
		i := h.FileInfo()
		if i.IsDir() {
			err := os.MkdirAll(p, i.Mode())
			if err != nil {
				return errors.Wrapf(err, "creating folder %q", p)
			}
			continue
		}
		err = func() error {
			fp, err := os.OpenFile(p, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, i.Mode())
			if err != nil {
				return errors.Wrapf(err, "extracting file %q", h.Name)
			}
			defer fp.Close()
			_, err = io.Copy(fp, tarFile)
			if err != nil {
				return errors.Wrapf(err, "writing file %q", h.Name)
			}
			return nil
		}()
		if err != nil {
			return errors.Wrap(err, "running extract")
		}
	}
	return nil
}

// InstallVersion downloads, extracts and installs the given go version.
func InstallVersion(v Version, src, targetPath string) error {
	response, err := http.Get(GODLURL + src)
	if err != nil {
		return errors.WithStack(err)
	}
	defer response.Body.Close()
	fmt.Println(GODLURL + "/" + src)
	gzFile, err := gzip.NewReader(response.Body)
	if err != nil {
		return errors.Wrap(err, "gunzipping file")
	}
	tarFile := tar.NewReader(gzFile)
	targetPath = filepath.Join(targetPath, v.String())
	if err := untar(tarFile, targetPath); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Update gets goVersion to the lates version of the
// given major.
func Update(goVersion string, updateRepos bool) {
}
