package cli

import (
	"context"
	"time"

	archives "deepsolutionsvn.com/disc/providers/musicbrainz/archives"
	commands "deepsolutionsvn.com/disc/providers/musicbrainz/commands"
)

func (opts *MusicBrainzPutArtistsJsonOptions) Execute(args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := ctrlcInit()

	stageCtx, err := prepareMusicBrainzStaging(archives.Artists, &standardOpts{
		dropPath:       opts.DropPath,
		dropId:         opts.DropId,
		libStagingPath: opts.LibStagingPath,
		quiet:          opts.Quiet,
	})
	if err != nil {
		return err
	}

	// fmt.Printf("%+v", stageCtx)

	err = commands.SaveJsonEntities(ctx, commands.SaveArtistsJson, stageCtx.archiveName, stageCtx.stagingBasePath, stageCtx.dropVersion, stageCtx.meter, interrupt)

	// if using progress meter, allow time for meter to update before exiting
	if stageCtx.meter != nil {
		time.Sleep(250 * time.Millisecond)
	}

	return err
}
