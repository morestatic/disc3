package commands

import (
	"encoding/json"
	"fmt"
	"os"

	archives "deepsolutionsvn.com/disc/discogs/archives"
	documents "deepsolutionsvn.com/disc/discogs/documents"
	indexes "deepsolutionsvn.com/disc/discogs/indexes"
)

func SaveDiscogsArtistJson(stagingPath string, indexer indexes.DiscogsArchiveIndexer, f *os.File, did int64) error {

	startPos, endPos, err := indexer.GetContentIdx(archives.Artists, did)
	if err != nil {
		return fmt.Errorf(fmt.Sprint("failed to read start and end pos: %w", err))
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

	filename, err := documents.WriteArtistInfoToStaging(stagingPath, did, string(prettyArtistJson))
	if err != nil {
		return fmt.Errorf("error writing artist json: %w", err)
	}
	fmt.Println(filename)

	artistReleases, err := indexer.GetArtistReleases(did)
	if err != nil {
		return fmt.Errorf("failed to get artist releases: %w", err)
	}

	artistReleasesByRole, err := indexes.GroupArtistReleasesByRole(artistReleases)
	if err != nil {
		return fmt.Errorf("error grouping artist releases: %w", err)
	}

	jsonReleases, err := json.MarshalIndent(artistReleasesByRole, "", "    ")
	if err != nil {
		return fmt.Errorf("error making json for artist releases: %w", err)
	}

	filename, err = documents.WriteArtistReleasesToStaging(stagingPath, did, string(jsonReleases))
	if err != nil {
		return fmt.Errorf("error writing releases json library staging: %w", err)
	}
	fmt.Println(filename)

	return nil
}
