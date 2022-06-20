package cli

import (
	"context"

	archives "deepsolutionsvn.com/disc/discogs/archives"
	commands "deepsolutionsvn.com/disc/discogs/commands"
)

func (opts *CreateMastersXmlIndexOptions) Execute(args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := ctrlcInit()

	idxCtx, err := prepareDiscogsIndex(ctx, archives.Masters, &standardOpts{
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

	err = commands.BuildMastersXmlIndex(ctx, idxCtx.d, idxCtx.i, idxCtx.m, done)
	return err
}
