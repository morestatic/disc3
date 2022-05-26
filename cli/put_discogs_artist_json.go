package cli

import (
	"fmt"
	"strconv"

	"deepsolutionsvn.com/disc/archives"
	"deepsolutionsvn.com/disc/commands"
	"deepsolutionsvn.com/disc/indexes"
)

func (opts *PutArtistJsonOptions) Execute(args []string) error {
	// fmt.Printf("%#v\n", *opts)

	did, err := strconv.ParseInt(opts.DiscId, 10, 64)
	if err != nil {
		return fmt.Errorf("unable to parse did: %w", err)
	}

	f, err := archives.NewFileReader(archives.Artists, opts.DropId, opts.DropPath)
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

	err = commands.SaveArtistJson(opts.LibStagingPath, indexer, f, did)
	if err != nil {
		return fmt.Errorf("error saving artist json: %w", err)
	}

	return nil
}
