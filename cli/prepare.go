package cli

import (
	"context"
	"fmt"
	"os"

	"deepsolutionsvn.com/disc/progress"
	discogs_archives "deepsolutionsvn.com/disc/providers/discogs/archives"
	indexes "deepsolutionsvn.com/disc/providers/discogs/indexes"
	musicbrainz_archives "deepsolutionsvn.com/disc/providers/musicbrainz/archives"
	musicbrainz_documents "deepsolutionsvn.com/disc/providers/musicbrainz/documents"
)

type standardOpts struct {
	dropId         string
	dropPath       string
	libStagingPath string
	connString     string
	storeType      string
	quiet          bool
}

type discogsIndexCtx struct {
	d *discogs_archives.DiscogsFileDecoder
	i indexes.DiscogsArchiveIndexer
	m progress.Meter
}

type discogsStageCtx struct {
	f *os.File
	i indexes.DiscogsArchiveIndexer
	m progress.Meter
}

type musicbrainzStageCtx struct {
	archiveName     string
	stagingBasePath string
	dropVersion     string
	meter           progress.Meter
}

func prepareDiscogsIndex(ctx context.Context, dt discogs_archives.DocumentType, opts *standardOpts) (*discogsIndexCtx, error) {

	d, err := discogs_archives.NewDiscogsFileDecoder(dt, opts.dropId, opts.dropPath)
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

func prepareDiscogsStaging(dt discogs_archives.DocumentType, opts *standardOpts) (*discogsStageCtx, error) {
	f, err := discogs_archives.NewDiscogsFileReader(dt, opts.dropId, opts.dropPath)
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

func prepareMusicBrainzStaging(dt musicbrainz_archives.DocumentType, opts *standardOpts) (stageCtx *musicbrainzStageCtx, err error) {
	archiveBasePath := opts.dropPath
	stagingBasePath := opts.libStagingPath
	dropVersion := opts.dropId

	if archiveBasePath == "" {
		archiveBasePath = os.Getenv("MUSICBRAINZ_DROP_PATH")
		if archiveBasePath == "" {
			archiveBasePath = musicbrainz_documents.MUSICBRAINZ_DROP_PATH
		}
	}

	if stagingBasePath == "" {
		stagingBasePath = os.Getenv("MUSICBRAINZ_STAGING_PATH")
		if stagingBasePath == "" {
			stagingBasePath = musicbrainz_documents.MUSICBRAINZ_STAGING_PATH
		}
	}

	if dropVersion == "" {
		dropVersion = os.Getenv("MUSICBRAINZ_DROP_VERSION")
		if dropVersion == "" {
			dropVersion = musicbrainz_documents.MUSICBRAINZ_DROP_VERSION
		}
	}

	var m progress.Meter = nil
	if !opts.quiet {
		m = progress.Setup(context.Background())
	}

	archiveName := musicbrainz_archives.MakeDefaultMusicBrainzArchiveName(dt, archiveBasePath, dropVersion)

	return &musicbrainzStageCtx{
			archiveName:     archiveName,
			stagingBasePath: stagingBasePath,
			dropVersion:     dropVersion,
			meter:           m},
		err
}
