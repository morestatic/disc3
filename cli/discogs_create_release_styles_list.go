package cli

import (
	"context"

	archives "deepsolutionsvn.com/disc/providers/discogs/archives"
	commands "deepsolutionsvn.com/disc/providers/discogs/commands"
)

func (opts *CreateReleaseStylesListOptions) Execute(args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := ctrlcInit()

	idxCtx, err := prepareDiscogsIndex(ctx, archives.Releases, &standardOpts{
		dropId:     opts.DropId,
		dropPath:   opts.DropPath,
		connString: opts.ConnString,
		storeType:  opts.StoreType,
		quiet:      opts.Quiet,
	})
	if err != nil {
		return err
	}
	defer idxCtx.d.Close()
	defer idxCtx.i.Close()

	err = commands.BuildReleaseStylesList(ctx, idxCtx.d, idxCtx.i, idxCtx.m, done)
	return err
}
