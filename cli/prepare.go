package cli

import (
	"context"
	"fmt"
	"os"

	"deepsolutionsvn.com/disc/progress"
	archives "deepsolutionsvn.com/disc/providers/discogs/archives"
	indexes "deepsolutionsvn.com/disc/providers/discogs/indexes"
)

type standardOpts struct {
	dropId     string
	dropPath   string
	connString string
	storeType  string
	quiet      bool
}

type discogsIndexCtx struct {
	d *archives.DiscogsFileDecoder
	i indexes.DiscogsArchiveIndexer
	m progress.Meter
}

type discogsStageCtx struct {
	f *os.File
	i indexes.DiscogsArchiveIndexer
	m progress.Meter
}

func prepareDiscogsIndex(ctx context.Context, dt archives.DocumentType, opts *standardOpts) (*discogsIndexCtx, error) {

	d, err := archives.NewDiscogsFileDecoder(dt, opts.dropId, opts.dropPath)
	if err != nil || d == nil {
		return nil, fmt.Errorf("failed to create new decoder: %w", err)
	}

	var i indexes.DiscogsArchiveIndexer
	if isDefaultIndexer(opts.storeType) {
		i, err = indexes.NewSQLiteDiscogsArchiveIndexer(opts.connString)
	} else {
		i, err = indexes.NewPGXArchiveIndexer(opts.connString)
	}
	if err != nil {
		d.Close()
		return nil, fmt.Errorf("unable to connect to indexer: %w", err)
	}

	var m progress.Meter = nil
	if !opts.quiet {
		m = progress.Setup(ctx)
	}

	return &discogsIndexCtx{d: d, i: i, m: m}, err
}

func prepareDiscogsStaging(dt archives.DocumentType, opts *standardOpts) (*discogsStageCtx, error) {
	f, err := archives.NewDiscogsFileReader(dt, opts.dropId, opts.dropPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read archive file: %w", err)
	}

	var i indexes.DiscogsArchiveIndexer
	if isDefaultIndexer(opts.storeType) {
		i, err = indexes.NewSQLiteDiscogsArchiveIndexer(opts.connString)
	} else {
		i, err = indexes.NewPGXArchiveIndexer(opts.connString)
	}
	if err != nil {
		f.Close()
		return nil, fmt.Errorf("unable to connect to indexer: %w", err)
	}

	var m progress.Meter = nil
	if !opts.quiet {
		m = progress.Setup(context.Background())
	}

	return &discogsStageCtx{f: f, i: i, m: m}, err
}
