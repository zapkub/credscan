package reposcan

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"

	"github.com/zapkub/credscan/internal"
	"github.com/zapkub/credscan/internal/appfs"
)

type RuleRunner func(ctx context.Context, filename string, reader io.Reader) ([]*internal.Finding, error)

type Runner struct {
	afs *appfs.AppFileSystem
}

func (rs *Runner) Scan(ctx context.Context, repositoryURL *url.URL, rules []RuleRunner) ([]*internal.Finding, error) {
	var err error
	var rt FileTraversal
	var findings = make([]*internal.Finding, 0)

	switch repositoryURL.Host {
	case "github.com":
		rt = NewGithubTraversal(repositoryURL, rs.afs)
	default:
		return nil, fmt.Errorf("unsupport repository %v", repositoryURL.Host)
	}

	err = rt.Clone()
	if err != nil {
		return nil, err
	}

	log.Println("start walk file....")
	rt.Walk(func(pathname string) error {
		log.Printf("scan file... %v", pathname)
		var err error
		for _, r := range rules {
			content, err := rt.Open(pathname)
			if err != nil {
				return err
			}
			filefindings, err := r(ctx, pathname, content)
			if err != nil {
				return err
			}
			findings = append(findings, filefindings...)
		}
		return err
	})
	log.Printf("walk end (%v)", findings)
	return findings, nil
}

func New(afs *appfs.AppFileSystem) *Runner {
	var rs Runner
	rs.afs = afs
	return &rs
}
