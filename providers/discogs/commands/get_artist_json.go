package commands

import (
	"fmt"
	"os"

	archives "deepsolutionsvn.com/disc/providers/discogs/archives"
	documents "deepsolutionsvn.com/disc/providers/discogs/documents"
	indexes "deepsolutionsvn.com/disc/providers/discogs/indexes"
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
