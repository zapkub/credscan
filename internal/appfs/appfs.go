package appfs

import (
	"io"
	"io/fs"
	"log"
	"os"
	"path"
)

type AppFileSystem struct {
	workingDir string
	appDir     string
	perm       fs.FileMode
}

func (afs *AppFileSystem) IsFileExistTemp(filename string) (bool, error) {
	p := path.Join(afs.workingDir, "temp", filename)
	if wdstate, err := os.Stat(p); os.IsNotExist(err) {
		return false, nil
	} else if wdstate.IsDir() {
		return false, nil
	} else if err != nil {
		return false, nil
	} else {
		return true, err
	}
}

func (appfs *AppFileSystem) OpenFromTemp(filename string, flag int) (*os.File, error) {
	return os.OpenFile(path.Join(appfs.workingDir, "/temp", filename), flag, appfs.perm)
}

func (appfs *AppFileSystem) SaveToTemp(reader io.Reader, outputname string) error {
	var outputpath = path.Join(appfs.workingDir, "/temp", outputname)
	err := os.MkdirAll(path.Dir(outputpath), appfs.perm)
	if err != nil {
		log.Println("unable to create output dir")
		return err
	}

	output, err := os.OpenFile(outputpath, os.O_WRONLY|os.O_CREATE, appfs.perm)
	if err != nil {
		log.Println("unable to open output file")
		return err
	}
	_, err = io.Copy(output, reader)
	return err
}

func (appfs *AppFileSystem) MustEnsureDir(pathname string) {
	mustEnsureDir(pathname, appfs.perm)
}
func (appfs *AppFileSystem) Pathname(pathname string) string {
	return path.Join(appfs.workingDir, pathname)
}

func (appfs *AppFileSystem) Unlink(pathname string) error {
	return os.RemoveAll(path.Join(appfs.workingDir, pathname))
}

func mustEnsureDir(dirname string, perm fs.FileMode) error {
	log.Printf("make sure %v is valid", dirname)
	if wdstate, err := os.Stat(dirname); !os.IsNotExist(err) && wdstate != nil {
		if !wdstate.IsDir() {
			log.Printf("%v is not a directory", dirname)
			os.Exit(1)
			return nil
		}
	} else if os.IsNotExist(err) {
		err = os.MkdirAll(dirname, perm)
		if err != nil {
			log.Printf("unable to create dir %v, %v", dirname, err)
			os.Exit(1)
		}
	} else {
		log.Printf("unexpected error, %v", err)
		panic(err)
	}
	return nil
}

type appFileSystemOptions struct {
	workingDir string
}
type appFileSystemOption func(*appFileSystemOptions)

func WorkingDir(dirname string) appFileSystemOption {
	return func(afso *appFileSystemOptions) {
		afso.workingDir = dirname
	}
}

func NewAppFileSystem(perm fs.FileMode, opts ...appFileSystemOption) *AppFileSystem {

	var afso appFileSystemOptions
	for _, opt := range opts {
		opt(&afso)
	}

	var appfs AppFileSystem
	appfs.perm = perm

	execdir, err := os.Executable()
	if err != nil {
		panic("unable to get executable dir")
	}
	appfs.appDir = path.Dir(execdir)

	if afso.workingDir == "" {
		homedir, err := os.UserHomeDir()
		if err != nil {
			panic("unable to get home dir")
		}
		appfs.workingDir = path.Join(homedir, ".credscan")
	} else {
		appfs.workingDir = afso.workingDir
	}
	mustEnsureDir(appfs.workingDir, appfs.perm)

	return &appfs
}
