package reposcan

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
)

type GithubFileTraversal struct {
	RepositoryURL *url.URL
	downloadURL   *url.URL
	name          string
	appfs         FileSystem
	httpclient    *http.Client
}

func (grt *GithubFileTraversal) getZipFilename() string {
	return grt.name + ".zip"
}

func (githubrt *GithubFileTraversal) Clone() error {
	var downloadurl = githubrt.downloadURL.String()
	resp, err := githubrt.httpclient.Get(downloadurl)
	if err != nil {
		log.Printf("unable to download file from %v", downloadurl)
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unable to download file status is %v, %v", resp.Status, downloadurl)
	}
	log.Printf("clone complete (%v)", downloadurl)
	return githubrt.appfs.SaveToTemp(resp.Body, githubrt.getZipFilename())
}

func (grt *GithubFileTraversal) IsClone() bool {
	isCloned, err := grt.appfs.IsFileExistTemp(grt.getZipFilename())
	if err != nil {
		log.Printf("unable to ensure zipfile or %v, %v", grt.name, err)
		return false
	}
	return isCloned
}

func (grt *GithubFileTraversal) Walk(wf WalkFunc) error {

	zipfile, err := grt.appfs.OpenFromTemp(grt.getZipFilename(), os.O_RDONLY)
	if err != nil {
		return err
	}
	fi, err := zipfile.Stat()
	if err != nil {
		return err
	}
	zipreader, err := zip.NewReader(zipfile, fi.Size())
	if err != nil {
		return err
	}
	for _, zf := range zipreader.File {
		if err := wf(zf.Name); err != nil {
			return err
		}
	}
	return nil
}

func (grt *GithubFileTraversal) Open(pathname string) (io.Reader, error) {
	zipfile, err := grt.appfs.OpenFromTemp(grt.getZipFilename(), os.O_RDONLY)
	if err != nil {
		return nil, err
	}
	fi, err := zipfile.Stat()
	if err != nil {
		return nil, err
	}
	zipreader, err := zip.NewReader(zipfile, fi.Size())
	if err != nil {
		return nil, err
	}
	for _, zf := range zipreader.File {
		if zf.Name == pathname {
			return zf.Open()
		}
	}
	return nil, fmt.Errorf("unable to open %v", pathname)
}

func NewGithubTraversal(u *url.URL, appfs FileSystem) *GithubFileTraversal {
	u.Host = "api.github.com"
	cloneURL, _ := url.Parse(u.String())
	cloneURL.Path = path.Join("repos", cloneURL.Path, "zipball", "master")

	var githubRi = GithubFileTraversal{
		RepositoryURL: u,
		downloadURL:   cloneURL,
		name:          u.Path,
		appfs:         appfs,
		httpclient:    http.DefaultClient,
	}

	return &githubRi
}
