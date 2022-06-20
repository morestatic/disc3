package commands

import (
	"context"
	"fmt"
	"io"

	discogstypes "deepsolutionsvn.com/disc/discogs/types"

	archives "deepsolutionsvn.com/disc/discogs/archives"
	indexes "deepsolutionsvn.com/disc/discogs/indexes"
	"deepsolutionsvn.com/disc/progress"
	"deepsolutionsvn.com/disc/utils"
)

func BuildArtistsXmlIndex(ctx context.Context, archiveDecoder *archives.DiscogsFileDecoder, indexer indexes.DiscogsArchiveIndexer, m progress.Meter, done chan struct{}) error {
	var ax *discogstypes.MinimalArtistXml = nil

	var count int64 = 1
	if m != nil {
		m.AddUnboundedBar("index", "indexing", func() string {
			return utils.Format(count)
		})
	}

	_, total, err := archiveDecoder.Scan(archives.Artists, "artist", discogstypes.MinimalArtistXml{}, func(d *archives.DiscogsFileDecoder, e interface{}, o interface{}, p archives.DiscogsScanProgress) (bool, error) {
		ax = o.(*discogstypes.MinimalArtistXml)

		err := indexer.AddWithContentIdx(archives.Artists, int64(ax.Id), p.StartPos, p.EndPos)
		if err != nil {
			return true, err
		}

		if m != nil {
			m.IncrUnboundedProgress("index", count)
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

	err = indexer.AddIdxCount(archives.Artists, total)
	if err != nil {
		return fmt.Errorf("error adding index count: %w", err)
	}

	if m != nil {
		m.SetTotal("index", 0, true)
	}

	return err
}
