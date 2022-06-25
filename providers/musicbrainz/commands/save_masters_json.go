package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"deepsolutionsvn.com/disc/providers/musicbrainz/archives"
	"deepsolutionsvn.com/disc/providers/musicbrainz/documents"
	"deepsolutionsvn.com/disc/providers/musicbrainz/scanner"
	"deepsolutionsvn.com/disc/utils"
)

type MasterInfo struct {
	Id string `json:"id"`
}

func SaveMastersJson(ctx context.Context, archivePath string, stagingBasePath string, dropVersion string, es *scanner.JSONEntityStream, interrupt chan struct{}, done chan struct{}, closer *sync.Once) {

	for entity := range es.Watch() {
		if entity.Error != nil {
			fmt.Printf("%s\n\n", entity.Error)
			break
		}

		prettyEntity, _ := json.MarshalIndent(entity.Value, "", "    ")

		var masterInfo MasterInfo
		_ = json.Unmarshal(prettyEntity, &masterInfo)

		_, _, _, _, err := documents.WriteFileToStaging(
			archives.Masters,
			stagingBasePath,
			dropVersion,
			masterInfo.Id,
			string(prettyEntity))

		if err != nil {
			fmt.Printf("%s\n\n", err)
			break
		}

		if utils.WasCancelled(ctx, interrupt, nil, done) {
			// make sure to exit if aborted
			break
		}
	}

	// ensure the the other goroutinues exit
	closer.Do(func() { close(done) })
}
