package commands

import (
	"fmt"
	"os"

	"deepsolutionsvn.com/disc/archives"
	"deepsolutionsvn.com/disc/documents"
	"deepsolutionsvn.com/disc/indexes"
)

func GetArtistJson(indexer indexes.DiscogsArchiveIndexer, f *os.File, did int64) error {
	startPos, endPos, err := indexer.GetContentIdx(archives.Artists, did)
	if err != nil {
		return fmt.Errorf("failed to read start and end pos: %w", err)
	}

	xmlContent, _, err := archives.GetContent(f, startPos, endPos)
	if err != nil {
		return fmt.Errorf("error reading xml content: %w", err)
	}

	json, err := documents.GetJson(xmlContent)
	if err != nil {
		return fmt.Errorf("error getting json: %w", err)
	}

	fmt.Println(json)
	return nil
}
