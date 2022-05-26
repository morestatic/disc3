package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"deepsolutionsvn.com/disc/archives"
	"deepsolutionsvn.com/disc/documents"
	"deepsolutionsvn.com/disc/indexes"
	"deepsolutionsvn.com/disc/progress"
)

func SaveLabelsJson(ctx context.Context, stagingPath string, indexer indexes.DiscogsArchiveIndexer, f *os.File, m progress.Meter, done chan struct{}) error {

	count := int64(1)
	total, err := indexer.GetIdxCount(indexes.IdxLabels)
	if err != nil {
		return fmt.Errorf("unable to get labels total: %w", err)
	}

	if m != nil {
		m.AddBar("stage", "staging", total)
	}

	for {
		start, end := indexes.GetBlockRange(count)

		// fmt.Printf("count = %d, start = %d, end = %d\n", count, start, end)

		labelsBatch, err := indexer.GetRangeOfDocumentIds(indexes.IdxLabels, start, end)
		if err != nil {
			return fmt.Errorf("error getting range of document ids: %w", err)
		}
		if len(labelsBatch) == 0 {
			break
		}

		for _, labelDid := range labelsBatch {

			if labelDid != 0 {

				exists, err := documents.LabelDocumentJsonExists(stagingPath, labelDid, "info")
				if err != nil {
					return fmt.Errorf("error checking if label document exists: %w", err)
				}

				if !exists {
					startPos, endPos, err := indexer.GetContentIdx(archives.Labels, labelDid)
					if err != nil {
						log.Printf("count = %d, did = %d\n", count, labelDid)
						return fmt.Errorf("failed to read start and end pos: %w", err)
					}

					xmlContent, _, err := archives.GetContent(f, startPos, endPos)
					if err != nil {
						return fmt.Errorf("error reading xml content: %w", err)
					}

					labelJsonStr, err := documents.GetJson(xmlContent)
					if err != nil {
						return fmt.Errorf("error getting json: %w", err)
					}

					var forPretty interface{}
					err = json.Unmarshal([]byte(labelJsonStr), &forPretty)
					if err != nil {
						return fmt.Errorf("error convert json string back to interface: %w", err)
					}

					prettyLabelJson, err := json.MarshalIndent(forPretty, "", "    ")
					if err != nil {
						return fmt.Errorf("error convert json string back to interface: %w", err)
					}

					// fmt.Println(string(prettyLabelJson))

					_, err = documents.WriteLabelInfoToStaging(stagingPath, labelDid, string(prettyLabelJson))
					if err != nil {
						return fmt.Errorf("error writing label json: %w", err)
					}
				}
			}

			if m != nil {
				m.IncrBar("stage", 1)
			}
			count++

			if wasCancelled(ctx, done) {
				return errCancelled
			}

		}
	}

	return nil
}
