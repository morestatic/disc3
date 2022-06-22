package commands

import (
	"context"
	"encoding/json"

	"deepsolutionsvn.com/disc/progress"
	"deepsolutionsvn.com/disc/providers/musicbrainz/archives"
	"deepsolutionsvn.com/disc/providers/musicbrainz/documents"
	"deepsolutionsvn.com/disc/providers/musicbrainz/scanner"
	"deepsolutionsvn.com/disc/utils"
)

type LabelInfo struct {
	Id string `json:"id"`
}

func SaveMusicBrainzLabelsJson(ctx context.Context, archivePath string, stagingBasePath string, dropVersion string, es *scanner.JSONEntityStream, m progress.Meter, done chan struct{}) error {

	// start the stream of json entities to save
	go es.Start(archivePath)

	count := int64(1)
	if m != nil {
		m.AddUnboundedBar("stage", "staging", func() string {
			return utils.Format(count)
		})
	}

	for entity := range es.Watch() {

		prettyEntity, _ := json.MarshalIndent(entity.Value, "", "    ")

		var labelInfo LabelInfo
		_ = json.Unmarshal(prettyEntity, &labelInfo)

		_, _, _, _, err := documents.WriteFileToStaging(
			archives.Labels,
			stagingBasePath,
			dropVersion,
			labelInfo.Id,
			string(prettyEntity))

		if err != nil {
			return err
		}

		if m != nil {
			m.IncrUnboundedProgress("stage", count)
		}
		count++

		if utils.WasCancelled(ctx, done) {
			return utils.ErrCancelled
		}

	}

	if m != nil {
		m.SetTotal("stage", 0, true)
	}

	return nil
}
