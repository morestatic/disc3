package commands

import (
	"fmt"
	"os"
	"strconv"

	archives "deepsolutionsvn.com/disc/discogs/archives"
	indexes "deepsolutionsvn.com/disc/discogs/indexes"
)

func GenerateArtistXml(archiveDecoder *archives.DiscogsFileDecoder, indexer indexes.DiscogsArchiveIndexer, f *os.File, w *archives.DiscogsFileEncoder, args []string) error {
	for _, arg := range args {
		did, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse artist id: %w", err)
		}

		startPos, endPos, err := indexer.GetContentIdx(archives.Artists, int64(did))
		if err != nil {
			return fmt.Errorf("failed to read start and end pos: %w", err)
		}

		xmlContent, _, err := archives.GetContent(f, startPos, endPos)
		if err != nil {
			return fmt.Errorf("error reading xml content: %w", err)
		}

		err = w.Write(xmlContent)
		if err != nil {
			return fmt.Errorf("error writing release xml content: %w", err)
		}

		fmt.Println(did)
	}

	return nil
}
