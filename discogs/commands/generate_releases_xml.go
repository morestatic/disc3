package commands

import (
	"fmt"
	"os"
	"strconv"

	archives "deepsolutionsvn.com/disc/discogs/archives"
	documents "deepsolutionsvn.com/disc/discogs/documents"
	indexes "deepsolutionsvn.com/disc/discogs/indexes"
)

func GenerateReleasesXml(stagingPath string, archiveDecoder *archives.DiscogsFileDecoder, indexer indexes.DiscogsArchiveIndexer, f *os.File, w *archives.DiscogsFileEncoder, args []string) error {
	existingReleases := make(map[documents.ReleaseDID]bool)
	for _, arg := range args {
		artistDid, err := strconv.ParseInt(arg, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse artist id: %w", err)
		}
		fmt.Println(artistDid)

		artistReleasesByRole, err := documents.ReadArtistReleasesJsonFromStaging(stagingPath, artistDid)
		if err != nil {
			return fmt.Errorf("failed to read artist releases from json library staging: %w", err)
		}

		for _, releasesByRole := range artistReleasesByRole {
			for _, releaseDid := range releasesByRole {
				if !existingReleases[releaseDid] {
					startPos, endPos, err := indexer.GetContentIdx(archives.Releases, int64(releaseDid))
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

					existingReleases[releaseDid] = true
				}
			}
		}

	}
	return nil
}
