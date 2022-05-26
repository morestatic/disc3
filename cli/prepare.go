package cli

import (
	"context"
	"fmt"
	"os"

	"deepsolutionsvn.com/disc/archives"
	"deepsolutionsvn.com/disc/indexes"
	"deepsolutionsvn.com/disc/progress"
)

type standardOpts struct {
	dropId     string
	dropPath   string
	connString string
	storeType  string
	quiet      bool
}

type indexCtx struct {
	d *archives.DiscogsFileDecoder
	i indexes.DiscogsArchiveIndexer
	m progress.Meter
}

type stageCtx struct {
	f *os.File
	i indexes.DiscogsArchiveIndexer
	m progress.Meter
}

func prepareDiscogsIndex(ctx context.Context, dt archives.DocumentType, opts *standardOpts) (*indexCtx, error) {

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

	return &indexCtx{d: d, i: i, m: m}, err
}

func prepareDiscogsStaging(dt archives.DocumentType, opts *standardOpts) (*stageCtx, error) {
	f, err := archives.NewFileReader(dt, opts.dropId, opts.dropPath)
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

	return &stageCtx{f: f, i: i, m: m}, err
}
