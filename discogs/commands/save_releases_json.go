package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	archives "deepsolutionsvn.com/disc/discogs/archives"
	documents "deepsolutionsvn.com/disc/discogs/documents"
	indexes "deepsolutionsvn.com/disc/discogs/indexes"
	"deepsolutionsvn.com/disc/progress"
)

func SaveDiscogsReleasesJson(ctx context.Context, stagingPath string, indexer indexes.DiscogsArchiveIndexer, f *os.File, m progress.Meter, done chan struct{}) error {

	count := int64(1)
	total, err := indexer.GetIdxCount(indexes.IdxReleases)
	if err != nil {
		return fmt.Errorf("unable to get releases total: %w", err)
	}

	if m != nil {
		m.AddBar("stage", "staging", total)
	}

	for {
		start, end := indexes.GetBlockRange(count)

		// fmt.Printf("count = %d, start = %d, end = %d\n", count, start, end)

		releasesBatch, err := indexer.GetRangeOfDocumentIds(indexes.IdxReleases, start, end)
		if err != nil {
			return fmt.Errorf("error getting range of document ids: %w", err)
		}
		if len(releasesBatch) == 0 {
			break
		}

		for _, releaseDid := range releasesBatch {

			if releaseDid != 0 {

				exists, err := documents.PrimaryDocumentExists(stagingPath, archives.Releases, releaseDid, "")
				if err != nil {
					return fmt.Errorf("error checking if release document exists: %w", err)
				}

				if !exists {
					startPos, endPos, err := indexer.GetContentIdx(archives.Releases, releaseDid)
					if err != nil {
						log.Printf("count = %d, did = %d\n", count, releaseDid)
						return fmt.Errorf("failed to read start and end pos: %w", err)
					}

					xmlContent, _, err := archives.GetContent(f, startPos, endPos)
					if err != nil {
						return fmt.Errorf("error reading xml content: %w", err)
					}

					releaseJsonStr, err := documents.GetJson(xmlContent)
					if err != nil {
						return fmt.Errorf("error getting json: %w", err)
					}

					var forPretty interface{}
					err = json.Unmarshal([]byte(releaseJsonStr), &forPretty)
					if err != nil {
						return fmt.Errorf("error convert json string back to interface: %w", err)
					}

					prettyReleaseJson, err := json.MarshalIndent(forPretty, "", "    ")
					if err != nil {
						return fmt.Errorf("error convert json string back to interface: %w", err)
					}

					// fmt.Println(string(prettyReleaseJson))

					_, _, _, _, err = documents.WriteFileToStaging(stagingPath, int64(releaseDid), string(prettyReleaseJson), archives.Releases, "")
					if err != nil {
						return fmt.Errorf("error writing release json library staging: %w", err)
					}
				}
			}

			if m != nil {
				m.IncrBar("stage", 1)
			}
			count++

			if WasCancelled(ctx, done) {
				return ErrCancelled
			}

		}
	}

	return nil
}
