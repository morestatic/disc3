package cli

import (
	"context"

	"deepsolutionsvn.com/disc/archives"
	"deepsolutionsvn.com/disc/commands"
)

func (opts *PutMastersJsonOptions) Execute(args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := ctrlcInit()

	stageCtx, err := prepareDiscogsStaging(archives.Masters, &standardOpts{
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

	err = commands.SaveMastersJson(ctx, opts.LibStagingPath, stageCtx.i, stageCtx.f, stageCtx.m, done)
	return err
}
