package commands

import (
	"encoding/json"
	"fmt"
	"os"

	archives "deepsolutionsvn.com/disc/providers/discogs/archives"
	documents "deepsolutionsvn.com/disc/providers/discogs/documents"
	indexes "deepsolutionsvn.com/disc/providers/discogs/indexes"
)

func SaveDiscogsArtistReleasesJson(stagingPath string, indexer indexes.DiscogsArchiveIndexer, f *os.File, did int64) error {

	artistReleasesByRole, err := documents.ReadArtistReleasesJsonFromStaging(stagingPath, did)
	if err != nil {
		return fmt.Errorf("failed to read artist releases from json library staging: %w", err)
	}

	for _, releasesByRole := range artistReleasesByRole {
		for _, releaseDid := range releasesByRole {

			startPos, endPos, err := indexer.GetContentIdx(archives.Releases, int64(releaseDid))
			if err != nil {
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

			filename, _, _, _, err := documents.WriteFileToStaging(stagingPath, int64(releaseDid), string(prettyReleaseJson), archives.Releases, "")
			if err != nil {
				return fmt.Errorf("error writing release json library staging: %w", err)
			}

			fmt.Println(filename)
		}
	}

	return nil
}
