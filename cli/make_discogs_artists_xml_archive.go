package cli

import (
	"errors"
	"fmt"

	"deepsolutionsvn.com/disc/archives"
	"deepsolutionsvn.com/disc/commands"
	"deepsolutionsvn.com/disc/indexes"
)

func (opts *MakeArtistsXmlArchiveOptions) Execute(args []string) error {
	fmt.Printf("%#v\n", *opts)
	fmt.Println(args)

	if len(args) == 0 {
		return errors.New("need at least one artist id")
	}

	archiveDecoder, err := archives.NewDiscogsFileDecoder(archives.Artists, opts.DropId, opts.DropPath)
	if err != nil || archiveDecoder == nil {
		return fmt.Errorf("failed to create new decoder: %w", err)
	}
	defer archiveDecoder.Close()

	f, err := archives.NewFileReader(archives.Artists, opts.DropId, opts.DropPath)
	if err != nil {
		return fmt.Errorf("unable to read archive file: %w", err)
	}
	defer f.Close()

	w, err := archives.NewDiscogsFileEncoder(archives.Artists, opts.OutputName)
	if err != nil {
		return fmt.Errorf("unable to create new file encoder: %w", err)
	}
	defer w.Close()

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

	err = commands.GenerateArtistXml(archiveDecoder, indexer, f, w, args)
	return err
}
