package update

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/kardianos/osext"

	"github.com/caicloud/nirvana/log"
	update "github.com/inconshreveable/go-update"
)

const (
	api = "https://api.github.com/repos/%v/releases/latest"
)

// GithubReleases ...
type GithubReleases struct {
	Asset []Asset `json:"assets"`
}

// Asset represents github release asset
type Asset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
	Version     string `json:"-"`
	Size        int64  `json:"size"`
}

// GetGithubLatestRelease returns latest release asset
func GetGithubLatestRelease(repo string) (Asset, error) {
	ret := Asset{}
	url := fmt.Sprintf(api, repo)
	resp, err := http.Get(url)
	if err != nil {
		return ret, err
	}

	defer resp.Body.Close()
	var m GithubReleases
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return ret, err
	}

	found := false
	for _, asset := range m.Asset {
		if strings.Contains(asset.Name, runtime.GOOS) {
			found = true
			ret = asset
			break
		}
	}
	if !found {
		return ret, fmt.Errorf("No release binary found on Github %v", url)
	}

	// https://github.com/caicloud/build-infra/releases/download/version/caimake-darwin
	slice := strings.Split(ret.DownloadURL, "/")
	if len(slice) != 9 {
		return ret, fmt.Errorf("Download url format error: %v", ret.DownloadURL)
	}

	ret.Version = slice[7]
	return ret, nil
}

// DoUpdate downloads binary from url and update itself
func DoUpdate(url string, size int64) error {

	// request the new file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// get executable path
	path, err := osext.Executable()
	if err != nil {
		return err
	}

	// apply update
	// Create our progress reporter and pass it to be used alongside our writer
	counter := &WriteCounter{
		Total: uint64(size),
	}
	log.Infof("Downloading file from %v", url)
	err = update.Apply(io.TeeReader(resp.Body, counter), update.Options{})
	if err != nil {
		if rerr := update.RollbackError(err); rerr != nil {
			log.Errorf("Failed to rollback from bad update: %v", rerr)
		}
	}

	log.Infof("Updated binary on %v", path)
	return err
}

// WriteCounter counts the number of bytes written to it. It implements to the io.Writer
// interface and we can pass this into io.TeeReader() which will report progress on each
// write cycle.
type WriteCounter struct {
	Total      uint64
	Downloaded uint64
	URL        string
	Reentry    bool
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Downloaded += uint64(n)
	wc.PrintProgress()
	wc.Reentry = true
	return n, nil
}

func (wc WriteCounter) PrintProgress() {
	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	if wc.Reentry {
		fmt.Print("\033[1A")
	}
	fmt.Printf("\rDownloading... %s/%s complete\n", humanize.Bytes(wc.Downloaded), humanize.Bytes(wc.Total))

}
