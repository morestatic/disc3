package cli

import (
	"context"
	"time"

	archives "deepsolutionsvn.com/disc/providers/musicbrainz/archives"
	commands "deepsolutionsvn.com/disc/providers/musicbrainz/commands"
	"deepsolutionsvn.com/disc/providers/musicbrainz/scanner"
)

func (opts *MusicBrainzPutMastersJsonOptions) Execute(args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := ctrlcInit()

	stageCtx, err := prepareMusicBrainzStaging(archives.Masters, &standardOpts{
		dropPath:       opts.DropPath,
		dropId:         opts.DropId,
		libStagingPath: opts.LibStagingPath,
		quiet:          opts.Quiet,
	})
	if err != nil {
		return err
	}

	es := scanner.NewJSONEntityStream()
	defer es.Close()

	err = commands.SaveMusicBrainzMastersJson(ctx, stageCtx.archiveName, stageCtx.stagingBasePath, stageCtx.dropVersion, es, stageCtx.meter, done)

	// if using progress meter, allow time for meter to update before exiting
	if stageCtx.meter != nil {
		time.Sleep(250 * time.Millisecond)
	}

	return err
}
