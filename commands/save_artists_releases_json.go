package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"deepsolutionsvn.com/disc/documents"
	"deepsolutionsvn.com/disc/indexes"
	"deepsolutionsvn.com/disc/progress"
)

func SaveArtistsReleasesJson(ctx context.Context, stagingPath string, indexer indexes.DiscogsArchiveIndexer, f *os.File, m progress.Meter, done chan struct{}) error {

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

		artistsBatch, err := indexer.GetRangeOfDocumentIds(indexes.IdxArtists, start, end)
		if err != nil {
			return fmt.Errorf("error getting range of document ids: %w", err)
		}
		if len(artistsBatch) == 0 {
			break
		}

		for _, artistDid := range artistsBatch {

			if artistDid != 0 {

				exists, err := documents.ArtistDocumentJsonExists(stagingPath, artistDid, "releases")
				if err != nil {
					return fmt.Errorf("error checking if artist document exists: %w", err)
				}

				if !exists {

					err = SaveArtistReleasesAsJsonDoc(artistDid, stagingPath, indexer)
					if err != nil {
						return err
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

func SaveArtistReleasesAsJsonDoc(artistDid int64, stagingPath string, indexer indexes.DiscogsArchiveIndexer) error {

	artistReleases, err := indexer.GetArtistReleases(artistDid)
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

	_, err = documents.WriteArtistReleasesToStaging(stagingPath, artistDid, string(jsonReleases))
	if err != nil {
		return fmt.Errorf("error writing artists releases json library staging: %w", err)
	}

	return nil
}
