package commands

import (
	"context"
	"fmt"
	"io"

	discogstypes "deepsolutionsvn.com/disc/types/discogs"

	"deepsolutionsvn.com/disc/archives"
	"deepsolutionsvn.com/disc/indexes"
	"deepsolutionsvn.com/disc/progress"
)

func BuildReleaseIndexes(ctx context.Context, archiveDecoder *archives.DiscogsFileDecoder, indexer indexes.DiscogsArchiveIndexer, m progress.Meter, done chan struct{}) error {
	var rx *discogstypes.ReleaseXmlWithArtistsAndLabels = nil

	total, err := indexer.GetIdxCount(indexes.IdxReleases)
	if err != nil {
		return fmt.Errorf("unable to get releases total: %w", err)
	}

	if m != nil {
		m.AddBar("index", "indexing", total)
	}

	_, _, err = archiveDecoder.Scan(archives.Releases, "release", discogstypes.ReleaseXmlWithArtistsAndLabels{}, func(d *archives.DiscogsFileDecoder, e interface{}, o interface{}, p archives.DiscogsScanProgress) (bool, error) {
		rx = o.(*discogstypes.ReleaseXmlWithArtistsAndLabels)

		tx, err := indexer.Begin()
		if err != nil {
			return true, fmt.Errorf("unable to start transaction: %w", err)
		}
		defer tx.Rollback()

		for _, label := range rx.Labels.Label {
			err = indexer.AddArtistRelease(label.Id, rx.Id, indexes.IdxPrimaryRole)
			if err != nil {
				return true, fmt.Errorf("unable to add artist release to label: %w", err)
			}
		}

		for _, artist := range rx.Artists.Artist {
			err = indexer.AddArtistRelease(artist.Id, rx.Id, indexes.IdxPrimaryRole)
			if err != nil {
				return true, fmt.Errorf("unable to add release for artist: %w", err)
			}
		}

		if len(rx.ExtraArtists.Artist) > 0 {
			for _, artist := range rx.ExtraArtists.Artist {
				err = indexer.AddArtistRelease(artist.Id, rx.Id, indexes.IdxExtraRole)
				if err != nil {
					return true, fmt.Errorf("unable to add release for extra artist: %w", err)
				}
			}
		}

		for _, track := range rx.TrackList.Tracks {
			if len(track.TrackArtists.Artist) > 0 {
				for _, artist := range track.TrackArtists.Artist {
					err = indexer.AddArtistRelease(artist.Id, rx.Id, indexes.IdxTrackRole)
					if err != nil {
						return true, fmt.Errorf("unable to add release for track artist: %w", err)
					}
				}
			}
		}

		err = tx.Commit()
		if err != nil {
			return true, fmt.Errorf("unable to commit releases transaction: %w", err)
		}

		if m != nil {
			m.IncrBar("index", 1)
		}

		if wasCancelled(ctx, done) {
			return true, errCancelled
		}

		return false, nil
	})

	if err == errCancelled {
		return err
	}

	if err != nil && err != io.EOF {
		return fmt.Errorf("error while scanning archive: %w", err)
	}

	return nil
}
