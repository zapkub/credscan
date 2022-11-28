package reposcan

import (
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/zapkub/credscan/internal/appfs"
	"github.com/zapkub/credscan/internal/testutil"
)

func newFixture() *GithubFileTraversal {
	wd, _ := os.Getwd()
	return NewGithubTraversal(
		testutil.MustParseURL("https://github.com/zapkub/react-thailand-address-typeahead"),
		appfs.NewAppFileSystem(
			0644,
			appfs.WorkingDir(path.Join(wd, "../../test/work_dir")),
		),
	)
}

func Test_NewGithubTraversal(t *testing.T) {
	githubrt := newFixture()
	if githubrt.name != "/zapkub/react-thailand-address-typeahead" {
		t.Errorf("repo name is not correct recieve %v", githubrt.name)
	}
}

func Test_Clone(t *testing.T) {

	githubrt := newFixture()
	var err = githubrt.Clone()
	if err != nil {
		t.Error(err)
	}

}

func Test_Walk(t *testing.T) {

	expect := []string{
		"zapkub-react-thailand-address-typeahead-849b865/",
		"zapkub-react-thailand-address-typeahead-849b865/.eslintrc",
		"zapkub-react-thailand-address-typeahead-849b865/.gitignore",
		"zapkub-react-thailand-address-typeahead-849b865/.npmignore",
		"zapkub-react-thailand-address-typeahead-849b865/.prettierrc",
		"zapkub-react-thailand-address-typeahead-849b865/.storybook/",
		"zapkub-react-thailand-address-typeahead-849b865/.storybook/main.js",
		"zapkub-react-thailand-address-typeahead-849b865/.storybook/preview.js",
		"zapkub-react-thailand-address-typeahead-849b865/README.md",
		"zapkub-react-thailand-address-typeahead-849b865/assets/",
		"zapkub-react-thailand-address-typeahead-849b865/assets/demo.gif",
		"zapkub-react-thailand-address-typeahead-849b865/jest.config.js",
		"zapkub-react-thailand-address-typeahead-849b865/lib/",
		"zapkub-react-thailand-address-typeahead-849b865/lib/context.ts",
		"zapkub-react-thailand-address-typeahead-849b865/lib/index.tsx",
		"zapkub-react-thailand-address-typeahead-849b865/lib/static-addr-source.ts",
		"zapkub-react-thailand-address-typeahead-849b865/lib/use-thailand-addr.test.ts",
		"zapkub-react-thailand-address-typeahead-849b865/lib/use-thailand-addr.ts",
		"zapkub-react-thailand-address-typeahead-849b865/package.json",
		"zapkub-react-thailand-address-typeahead-849b865/stories/",
		"zapkub-react-thailand-address-typeahead-849b865/stories/Page.css",
		"zapkub-react-thailand-address-typeahead-849b865/stories/Page.stories.tsx",
		"zapkub-react-thailand-address-typeahead-849b865/tsconfig.json",
		"zapkub-react-thailand-address-typeahead-849b865/webpack.config.js",
		"zapkub-react-thailand-address-typeahead-849b865/yarn.lock",
	}
	var pathnames []string

	githubrt := newFixture()
	if !githubrt.IsClone() {
		t.Errorf("file should be existed") // checking on /test/work_dir
	}

	var err = githubrt.Walk(func(pathname string) error {
		pathnames = append(pathnames, pathname)
		return nil
	})
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expect, pathnames) {
		t.Error("result does not match expect, walk should traverse to every folder and file")
	}

}
