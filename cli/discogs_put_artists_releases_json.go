package cli

import (
	"context"

	archives "deepsolutionsvn.com/disc/discogs/archives"
	commands "deepsolutionsvn.com/disc/discogs/commands"
)

func (opts *DiscogsPutArtistsReleasesJsonOptions) Execute(args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := ctrlcInit()

	stageCtx, err := prepareDiscogsStaging(archives.Artists, &standardOpts{
		dropId:     opts.DropId,
		dropPath:   opts.DropPath,
		connString: opts.ConnString,
		storeType:  opts.StoreType,
		quiet:      opts.Quiet,
	})
	if err != nil {
		return err
	}
	defer stageCtx.f.Close()
	defer stageCtx.i.Close()

	err = commands.SaveDiscogsArtistsReleasesJson(ctx, opts.LibStagingPath, stageCtx.i, stageCtx.f, stageCtx.m, done)
	return err

}
