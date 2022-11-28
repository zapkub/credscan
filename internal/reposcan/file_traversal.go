package reposcan

import (
	"io"
	"net/url"
	"os"
)

type FileSystem interface {
	IsFileExistTemp(filename string) (bool, error)
	OpenFromTemp(pathname string, flag int) (*os.File, error)
	SaveToTemp(src io.Reader, outname string) error
}

type TextFile struct {
	ContentURL *url.URL
}

type WalkFunc func(pathname string) error

// FileTraversal (repository traversal) is a simple interface
// for iterate inside repository
type FileTraversal interface {
	Clone() error
	Walk(WalkFunc) error
	Open(pathname string) (io.Reader, error)
}
