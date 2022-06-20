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

func BuildReleasesXmlIndex(ctx context.Context, archiveDecoder *archives.DiscogsFileDecoder, indexer indexes.DiscogsArchiveIndexer, m progress.Meter, done chan struct{}) error {
	var rx *discogstypes.MinimalReleaseXml = nil

	var count int64 = 1
	if m != nil {
		m.AddUnboundedBar("index", "indexing", func() string {
			return utils.Format(count)
		})
	}

	_, total, err := archiveDecoder.Scan(archives.Releases, "release", discogstypes.MinimalReleaseXml{}, func(d *archives.DiscogsFileDecoder, e interface{}, o interface{}, p archives.DiscogsScanProgress) (bool, error) {
		rx = o.(*discogstypes.MinimalReleaseXml)

		err := indexer.AddWithContentIdx(archives.Releases, int64(rx.Id), p.StartPos, p.EndPos)
		if err != nil {
			return true, fmt.Errorf("unable to add content with index (%d): %w", rx.Id, err)
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

	err = indexer.AddIdxCount(archives.Releases, total)
	if err != nil {
		return fmt.Errorf("error while index count: %w", err)
	}

	if m != nil {
		m.SetTotal("index", 0, true)
	}

	return err
}
