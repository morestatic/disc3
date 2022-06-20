package cli

import (
	"fmt"
	"strconv"

	archives "deepsolutionsvn.com/disc/providers/discogs/archives"
	commands "deepsolutionsvn.com/disc/providers/discogs/commands"
	indexes "deepsolutionsvn.com/disc/providers/discogs/indexes"
)

func (opts *ExtractReleaseJsonOptions) Execute(args []string) error {
	// fmt.Printf("%#v\n", *c)

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
		return fmt.Errorf("unable to connect to db: %w", err)
	}
	defer indexer.Close()

	err = commands.GetReleaseJson(indexer, f, did)
	return err
}
