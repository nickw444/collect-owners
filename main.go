package main

import (
	"os"

	"github.com/nickw444/collect-owners/usernamedb"
	"gopkg.in/alecthomas/kingpin.v2"
	"regexp"
)

var Version string

func main() {

	var (
		app = kingpin.New("Collect Owners", "Walk A Repo and compile a Github CODEOWNERS file").Version(Version)

		repo          = app.Arg("repo", "Path to repository").Required().String()
		contributors  = app.Flag("contributors", "Path to contributors file to add to the users DB").String()
		excludesRaw   = app.Flag("exclude", "Regular expressions of paths to exclude").Strings()
		addUnresolved = app.Flag("add-unresolved", "Add ownerships that do not have entries in the users DB as their raw entries within the OWNERS files").Bool()
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	var err error
	repoRoot := *repo

	var dbLoader usernamedb.DBLoader
	if *contributors != "" {
		dbLoader = &usernamedb.ContributorsFileDBLoader{
			Filename: *contributors,
		}
	} else {
		dbLoader = &usernamedb.SimpleDBLoader{}
	}

	usernameDb := &usernamedb.UsernameDB{
		Loader:        dbLoader,
		AddUnresolved: *addUnresolved,
	}
	err = usernameDb.Load()
	if err != nil {
		panic(err)
	}

	formatter := &OwnersFileFormatter{
		UsernameDB: usernameDb,
	}

	fileProcessor := &OwnersFileProcessor{
		RepoRoot: repoRoot,
	}

	var excludes []*regexp.Regexp
	for _, excludeRaw := range *excludesRaw {
		excludes = append(excludes, regexp.MustCompile(excludeRaw))
	}

	treeProvider := OwnersTreeProvider{
		ownersFileProcessor: fileProcessor,
		RootPath:            repoRoot,
		Excludes:            excludes,
	}

	rootEnt, err := treeProvider.GetFileTree()
	if err != nil {
		panic(err)
	}

	formatter.Format(rootEnt)
}
