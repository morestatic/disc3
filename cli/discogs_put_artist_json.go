package cli

import (
	"fmt"
	"strconv"

	archives "deepsolutionsvn.com/disc/discogs/archives"
	commands "deepsolutionsvn.com/disc/discogs/commands"
	indexes "deepsolutionsvn.com/disc/discogs/indexes"
)

func (opts *DiscogsPutArtistJsonOptions) Execute(args []string) error {
	// fmt.Printf("%#v\n", *opts)

	did, err := strconv.ParseInt(opts.DiscId, 10, 64)
	if err != nil {
		return fmt.Errorf("unable to parse did: %w", err)
	}

	f, err := archives.NewDiscogsFileReader(archives.Artists, opts.DropId, opts.DropPath)
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
		return fmt.Errorf("unable to connect to db:%w", err)
	}
	defer indexer.Close()

	err = commands.SaveDiscogsArtistJson(opts.LibStagingPath, indexer, f, did)
	if err != nil {
		return fmt.Errorf("error saving artist json: %w", err)
	}

	return nil
}
