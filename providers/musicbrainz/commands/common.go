package commands

import (
	"context"
	"sync"
	"time"

	"deepsolutionsvn.com/disc/progress"
	"deepsolutionsvn.com/disc/providers/musicbrainz/scanner"
	"deepsolutionsvn.com/disc/utils"
)

type SlurpFunc func(ctx context.Context, archivePath string, stagingBasePath string, dropVersion string, es *scanner.JSONEntityStream, interrupt chan struct{}, done chan struct{}, closer *sync.Once)

func SaveJsonEntities(ctx context.Context, slurpFunc SlurpFunc, archivePath string, stagingBasePath string, dropVersion string, m progress.Meter, interrupt chan struct{}) error {
	// guard against potential race conditions when closing the done channel
	// see https://groups.google.com/g/golang-nuts/c/rhxMiNmRAPk
	closer := &sync.Once{}

	done := make(chan struct{})

	var wg sync.WaitGroup

	// start the stream of json entities to be staged
	es := scanner.NewJSONEntityStream()
	go es.Start(archivePath)

	if m != nil {
		initProgress(es, m)
	}

	if m != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			updateProgress(ctx, es, m, interrupt, done)
		}()
	}

	for i := 1; i < scanner.EntityChannelBufferSize; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			slurpFunc(ctx, archivePath, stagingBasePath, dropVersion, es, interrupt, done, closer)
		}()
	}

	wg.Wait()

	if m != nil {
		m.SetTotal("stage", 0, true)
	}

	return nil
}

func initProgress(es *scanner.JSONEntityStream, m progress.Meter) {
	m.AddUnboundedBar("stage", "staging", func() string {
		return utils.Format(es.GetCount())
	})
	m.IncrUnboundedProgress("stage", 1)
}

func updateProgress(ctx context.Context, es *scanner.JSONEntityStream, m progress.Meter, interrupt chan struct{}, done chan struct{}) {
	for {
		time.Sleep(250 * time.Millisecond)

		m.SetBar("stage", es.GetCount())

		if utils.WasCancelled(ctx, interrupt, nil, done) {
			// finish if interrupted
			return
		}
	}
}
