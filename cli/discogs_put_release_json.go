package cli

import (
	"fmt"
	"log"
	"strconv"

	archives "deepsolutionsvn.com/disc/providers/discogs/archives"
	commands "deepsolutionsvn.com/disc/providers/discogs/commands"
	indexes "deepsolutionsvn.com/disc/providers/discogs/indexes"
)

func (opts *DiscogsPutReleaseJsonOptions) Execute(args []string) error {
	// fmt.Printf("%#v\n", *opts)

	did, err := strconv.ParseInt(opts.DiscId, 10, 64)
	if err != nil {
		return err
	}

	f, err := archives.NewDiscogsFileReader(archives.Releases, opts.DropId, opts.DropPath)
	if err != nil {
		return fmt.Errorf("unable to read archive file: %w", err)
	}
	defer f.Close()

	var indexer indexes.DiscogsArchiveIndexer
	if isDefaultIndexer(opts.StoreType) {
		indexer, err = indexes.NewSQLiteDiscogsArchiveIndexer(opts.ConnString)
	} else {
		indexer, err = indexes.NewPGXArchiveIndexer(opts.ConnString)
	}
	if err != nil {
		log.Fatal("unable to connect to db: ", err)
	}
	defer indexer.Close()

	err = commands.SaveDiscogsReleaseJson(opts.LibStagingPath, indexer, f, did)
	return err
}
