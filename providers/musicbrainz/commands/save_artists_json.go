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

type ArtistInfo struct {
	Id string `json:"id"`
}

func SaveArtistsJson(ctx context.Context, archivePath string, stagingBasePath string, dropVersion string, es *scanner.JSONEntityStream, interrupt chan struct{}, done chan struct{}, closer *sync.Once) {

	for entity := range es.Watch() {
		if entity.Error != nil {
			fmt.Printf("%s\n\n", entity.Error)
			break
		}

		prettyEntity, _ := json.MarshalIndent(entity.Value, "", "    ")

		var artistInfo ArtistInfo
		_ = json.Unmarshal(prettyEntity, &artistInfo)

		_, _, _, _, err := documents.WriteFileToStaging(
			archives.Artists,
			stagingBasePath,
			dropVersion,
			artistInfo.Id,
			string(prettyEntity))

		if err != nil {
			fmt.Printf("%s\n\n", err)
			break
		}

		if utils.WasCancelled(ctx, interrupt, nil, done) {
			break
		}
	}

	// ensure the the other goroutinues exit
	closer.Do(func() { close(done) })
}
