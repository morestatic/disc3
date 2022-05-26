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

func SaveArtistsJson(ctx context.Context, stagingPath string, indexer indexes.DiscogsArchiveIndexer, f *os.File, m progress.Meter, done chan struct{}) error {

	count := int64(1)
	total, err := indexer.GetIdxCount(indexes.IdxArtists)
	if err != nil {
		return fmt.Errorf("unable to get artists total: %w", err)
	}

	if m != nil {
		m.AddBar("stage", "staging", total)
	}

	for {
		start, end := indexes.GetBlockRange(count)

		// fmt.Printf("count = %d, start = %d, end = %d\n", count, start, end)

		artistsBatch, err := indexer.GetRangeOfDocumentIds(indexes.IdxArtists, start, end)
		if err != nil {
			return fmt.Errorf("error getting range of document ids: %w", err)
		}
		if len(artistsBatch) == 0 {
			break
		}

		for _, artistDid := range artistsBatch {

			if artistDid != 0 {

				exists, err := documents.ArtistDocumentJsonExists(stagingPath, artistDid, "info")
				if err != nil {
					return fmt.Errorf("error checking if artist document exists: %w", err)
				}

				if !exists {
					startPos, endPos, err := indexer.GetContentIdx(archives.Artists, artistDid)
					if err != nil {
						log.Printf("count = %d, did = %d\n", count, artistDid)
						return fmt.Errorf("failed to read start and end pos: %w", err)
					}

					xmlContent, _, err := archives.GetContent(f, startPos, endPos)
					if err != nil {
						return fmt.Errorf("error reading xml content: %w", err)
					}

					artistJsonStr, err := documents.GetJson(xmlContent)
					if err != nil {
						return fmt.Errorf("error getting json: %w", err)
					}

					var forPretty interface{}
					err = json.Unmarshal([]byte(artistJsonStr), &forPretty)
					if err != nil {
						return fmt.Errorf("error convert json string back to interface: %w", err)
					}

					prettyArtistJson, err := json.MarshalIndent(forPretty, "", "    ")
					if err != nil {
						return fmt.Errorf("error convert json string back to interface: %w", err)
					}

					// fmt.Println(string(prettyArtistJson))

					_, err = documents.WriteArtistInfoToStaging(stagingPath, artistDid, string(prettyArtistJson))
					if err != nil {
						return fmt.Errorf("error writing artist json: %w", err)
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
