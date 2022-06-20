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

func BuildLabelsXmlIndex(ctx context.Context, archiveDecoder *archives.DiscogsFileDecoder, indexer indexes.DiscogsArchiveIndexer, m progress.Meter, done chan struct{}) error {
	var lx *discogstypes.MinimalLabelXml = nil
	var count int64 = 1

	if m != nil {
		m.AddUnboundedBar("index", "indexing", func() string {
			return utils.Format(count)
		})
	}

	tx, err := indexer.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Rollback()

	_, total, err := archiveDecoder.Scan(archives.Labels, "label", discogstypes.MinimalLabelXml{}, func(d *archives.DiscogsFileDecoder, e interface{}, o interface{}, p archives.DiscogsScanProgress) (bool, error) {
		lx = o.(*discogstypes.MinimalLabelXml)

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

		err = indexer.AddWithContentIdx(archives.Labels, int64(lx.Id), p.StartPos, p.EndPos)
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

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	if err == ErrCancelled {
		return err
	}

	if err != nil && err != io.EOF {
		return fmt.Errorf("error while scanning archive: %w", err)
	}

	err = indexer.AddIdxCount(archives.Labels, total)
	if err != nil {
		return fmt.Errorf("error adding index count: %w", err)
	}

	if m != nil {
		m.SetTotal("index", 0, true)
	}

	return err
}
