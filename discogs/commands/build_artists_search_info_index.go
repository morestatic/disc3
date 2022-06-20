package commands

import (
	"context"
	"fmt"
	"io"

	archives "deepsolutionsvn.com/disc/discogs/archives"
	indexes "deepsolutionsvn.com/disc/discogs/indexes"
	discogstypes "deepsolutionsvn.com/disc/discogs/types"
	"deepsolutionsvn.com/disc/progress"
)

func BuildArtistsSearchInfo(ctx context.Context, archiveDecoder *archives.DiscogsFileDecoder, indexer indexes.DiscogsArchiveIndexer, m progress.Meter, done chan struct{}) error {
	var ax *discogstypes.ArtistXmlWithInfo = nil
	var count int64 = 1

	total, err := indexer.GetIdxCount(indexes.IdxArtists)
	if err != nil {
		return fmt.Errorf("unable to get artists total: %w", err)
	}

	if m != nil {
		m.AddBar("index", "indexing", total)
	}

	tx, err := indexer.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Rollback()

	_, _, err = archiveDecoder.Scan(archives.Artists, "artist", discogstypes.ArtistXmlWithInfo{}, func(d *archives.DiscogsFileDecoder, e interface{}, o interface{}, p archives.DiscogsScanProgress) (bool, error) {
		ax = o.(*discogstypes.ArtistXmlWithInfo)

		as := &discogstypes.ArtistSearchInfo{
			Id:       ax.Id,
			Name:     ax.Name,
			RealName: ax.RealName,
			IsGroup:  len(ax.MemberIds) > 0,
		}

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

		err := indexer.AddArtistSearchInfo(int64(ax.Id), as)
		if err != nil {
			return true, err
		}

		if m != nil {
			m.IncrBar("index", 1)
		}
		count++

		if WasCancelled(ctx, done) {
			return true, ErrCancelled
		}

		return false, nil
	})

	if err == ErrCancelled {
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
