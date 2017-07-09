package main

import (
	"os"

	"github.com/nickw444/collect-owners/usernamedb"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {

	var (
		app = kingpin.New("Collect Owners", "Walk A Repo and compile a Github CODEOWNERS file")

		repo          = app.Arg("repo", "Path to repository").Required().String()
		contributors  = app.Flag("contributors", "Path to contributors file to add to the users DB").String()
		addUnresolved = app.Flag("add-unresolved", "Add ownerships that do not have entries in the users DB as their raw entries within the OWNERS files").Bool()
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	var err error
	repoRoot := *repo

	formatter := &OwnersFileFormatter{}
	fileProvider := &FileProvider{
		RepoRoot: repoRoot,
	}
	fileProcessor := &OwnersFileProcessor{
		RepoRoot: repoRoot,
	}

	var dbLoader usernamedb.DBLoader
	if contributors != nil {
		dbLoader = &usernamedb.ContributorsFileDBLoader{
			Filename: *contributors,
		}
	} else {
		dbLoader = &usernamedb.SimpleDBLoader{}
	}

	usernameDb := usernamedb.UsernameDB{
		Loader:        dbLoader,
		AddUnresolved: *addUnresolved,
	}
	err = usernameDb.Load()
	if err != nil {
		panic(err)
	}

	walker := OwnersWalker{
		UsernameDB:    usernameDb,
		RepoRoot:      repoRoot,
		FileProvider:  fileProvider,
		FileProcessor: fileProcessor,
	}

	err = walker.Walk()
	if err != nil {
		panic(err)
	}
	entries := walker.CollectEntries()
	formatter.Format(entries)
}
