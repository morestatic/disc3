package cli

import (
	"context"
	"time"

	archives "deepsolutionsvn.com/disc/providers/musicbrainz/archives"
	commands "deepsolutionsvn.com/disc/providers/musicbrainz/commands"
)

func (opts *MusicBrainzPutLabelsJsonOptions) Execute(args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := ctrlcInit()

	stageCtx, err := prepareMusicBrainzStaging(archives.Labels, &standardOpts{
		dropPath:       opts.DropPath,
		dropId:         opts.DropId,
		libStagingPath: opts.LibStagingPath,
		quiet:          opts.Quiet,
	})
	if err != nil {
		return err
	}

	err = commands.SaveJsonEntities(ctx, commands.SaveLabelsJson, stageCtx.archiveName, stageCtx.stagingBasePath, stageCtx.dropVersion, stageCtx.meter, interrupt)

	// if using progress meter, allow time for meter to update before exiting
	if stageCtx.meter != nil {
		time.Sleep(250 * time.Millisecond)
	}

	return err
}
