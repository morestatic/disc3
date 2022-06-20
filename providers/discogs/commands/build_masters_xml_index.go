package commands

import (
	"context"
	"fmt"
	"io"

	discogstypes "deepsolutionsvn.com/disc/providers/discogs/types"

	"deepsolutionsvn.com/disc/progress"
	archives "deepsolutionsvn.com/disc/providers/discogs/archives"
	indexes "deepsolutionsvn.com/disc/providers/discogs/indexes"
	"deepsolutionsvn.com/disc/utils"
)

func BuildMastersXmlIndex(ctx context.Context, archiveDecoder *archives.DiscogsFileDecoder, indexer indexes.DiscogsArchiveIndexer, m progress.Meter, done chan struct{}) error {
	var mx *discogstypes.MinimalMasterXml = nil

	var count int64 = 1
	if m != nil {
		m.AddUnboundedBar("index", "indexing", func() string {
			return utils.Format(count)
		})
	}

	_, total, err := archiveDecoder.Scan(archives.Masters, "master", discogstypes.MinimalMasterXml{}, func(d *archives.DiscogsFileDecoder, e interface{}, o interface{}, p archives.DiscogsScanProgress) (bool, error) {
		mx = o.(*discogstypes.MinimalMasterXml)

		err := indexer.AddWithContentIdx(archives.Masters, int64(mx.Id), p.StartPos, p.EndPos)
		if err != nil {
			return true, fmt.Errorf("unable to add content with index (%d): %w", mx.Id, err)
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

	err = indexer.AddIdxCount(archives.Masters, total)
	if err != nil {
		return fmt.Errorf("error adding index count: %w", err)
	}

	if m != nil {
		m.SetTotal("index", 0, true)
	}

	return err
}
