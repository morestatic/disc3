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

func SaveMastersJson(ctx context.Context, stagingPath string, indexer indexes.DiscogsArchiveIndexer, f *os.File, m progress.Meter, done chan struct{}) error {

	count := int64(1)
	total, err := indexer.GetIdxCount(indexes.IdxMasters)
	if err != nil {
		return fmt.Errorf("unable to get masters total: %w", err)
	}

	if m != nil {
		m.AddBar("stage", "staging", total)
	}

	for {
		start, end := indexes.GetBlockRange(count)

		// fmt.Printf("count = %d, start = %d, end = %d\n", count, start, end)

		mastersBatch, err := indexer.GetRangeOfDocumentIds(indexes.IdxMasters, start, end)
		if err != nil {
			return fmt.Errorf("error getting range of document ids: %w", err)
		}
		if len(mastersBatch) == 0 {
			break
		}

		for _, masterDid := range mastersBatch {

			if masterDid != 0 {

				exists, err := documents.PrimaryDocumentExists(stagingPath, archives.Masters, masterDid, "")
				if err != nil {
					return fmt.Errorf("error checking if master document exists: %w", err)
				}

				if !exists {
					startPos, endPos, err := indexer.GetContentIdx(archives.Masters, masterDid)
					if err != nil {
						log.Printf("count = %d, did = %d\n", count, masterDid)
						return fmt.Errorf("failed to read start and end pos: %w", err)
					}

					xmlContent, _, err := archives.GetContent(f, startPos, endPos)
					if err != nil {
						log.Printf("did = %d, start_pos = %d, end_pos = %d\n", masterDid, startPos, endPos)
						return fmt.Errorf("error reading xml content: %w", err)
					}

					masterJsonStr, err := documents.GetJson(xmlContent)
					if err != nil {
						return fmt.Errorf("error getting json: %w", err)
					}

					var forPretty interface{}
					err = json.Unmarshal([]byte(masterJsonStr), &forPretty)
					if err != nil {
						return fmt.Errorf("error convert json string back to interface: %w", err)
					}

					prettyMasterJson, err := json.MarshalIndent(forPretty, "", "    ")
					if err != nil {
						return fmt.Errorf("error convert json string back to interface: %w", err)
					}

					// fmt.Println(string(prettyMasterJson))

					_, _, _, _, err = documents.WriteFileToStaging(stagingPath, int64(masterDid), string(prettyMasterJson), archives.Masters, "")
					if err != nil {
						return fmt.Errorf("error writing master json library staging: %w", err)
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
