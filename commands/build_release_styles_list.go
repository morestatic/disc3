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

func BuildReleaseStylesList(ctx context.Context, archiveDecoder *archives.DiscogsFileDecoder, indexer indexes.DiscogsArchiveIndexer, m progress.Meter, done chan struct{}) error {
	var rx *discogstypes.ReleaseXmlWithStyles = nil

	var count int64 = 1
	total, err := indexer.GetIdxCount(indexes.IdxReleases)
	if err != nil {
		return fmt.Errorf("unable to get releases total: %w", err)
	}

	if m != nil {
		m.AddBar("index", "indexing", total)
	}

	tx, err := indexer.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Rollback()

	_, _, err = archiveDecoder.Scan(archives.Releases, "release", discogstypes.ReleaseXmlWithStyles{}, func(d *archives.DiscogsFileDecoder, e interface{}, o interface{}, p archives.DiscogsScanProgress) (bool, error) {
		rx = o.(*discogstypes.ReleaseXmlWithStyles)

		if count%1000 == 0 {
			err = tx.Commit()
			if err != nil {
				return true, fmt.Errorf("unable to commit transaction: %w", err)
			}

			tx, err = indexer.Begin()
			if err != nil {
				return true, fmt.Errorf("unable to start transaction: %w", err)
			}
		}

		for _, s := range rx.StylesList.Styles {
			rs := &discogstypes.ReleaseStyleEntry{
				Style: s,
			}
			err := indexer.AddStyle(rs)
			if err != nil {
				return true, err
			}
		}

		if m != nil {
			m.IncrBar("index", 1)
		}
		count++

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

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil
}
