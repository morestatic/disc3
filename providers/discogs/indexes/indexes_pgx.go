package indexes

import (
	"context"
	"errors"
	"fmt"

	discogstypes "deepsolutionsvn.com/disc/providers/discogs/types"

	archives "deepsolutionsvn.com/disc/providers/discogs/archives"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PGXArchiveIndexer struct {
	dbp *pgxpool.Pool
}

type PGXIndexerTx struct {
	tx pgx.Tx
	DiscogsArchiveIndexerTx
}

func NewPGXArchiveIndexer(connString string) (*PGXArchiveIndexer, error) {
	if connString == "" {
		connString = GetDefaultConnUrl()
	}

	dbpool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	connector := &PGXArchiveIndexer{
		dbp: dbpool,
	}

	return connector, nil
}

func newPGXIndexerTx(newTx pgx.Tx) *PGXIndexerTx {
	ptx := &PGXIndexerTx{
		tx: newTx,
	}
	return ptx
}

func (p *PGXArchiveIndexer) Close() {
	p.dbp.Close()
}

func (p *PGXArchiveIndexer) GetContentIdx(dt archives.DocumentType, did int64) (int64, int64, error) {
	var startPos, endPos int64

	sql := fmt.Sprintf("SELECT start_pos, end_pos FROM %s WHERE did = $1", dt)
	err := p.dbp.QueryRow(context.Background(), sql, did).Scan(
		&startPos,
		&endPos,
	)
	if err != nil {
		return 0, 0, err
	}

	return startPos, endPos, nil
}

func (p *PGXArchiveIndexer) AddWithContentIdx(dt archives.DocumentType, did int64, startPos int64, endPos int64) error {
	sql := fmt.Sprintf("INSERT INTO %s (did, start_pos, end_pos) VALUES ($1, $2, $3)", dt)

	commandTag, err := p.dbp.Exec(context.Background(), sql, did, startPos, endPos)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("item not added")
	}

	return nil
}

func (p *PGXArchiveIndexer) UpdateContentIdx(dt archives.DocumentType, did int64, startPos int64, endPos int64) error {
	sql := fmt.Sprintf("UPDATE %s SET start_pos = $2, end_pos = $3 WHERE did = $1", dt)

	commandTag, err := p.dbp.Exec(context.Background(), sql, did, startPos, endPos)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("content idx not updated")
	}

	return nil
}

func (p *PGXArchiveIndexer) AddIdxCount(dt archives.DocumentType, count int64) error {
	return errors.New("not yet implemented")
}

func (p *PGXArchiveIndexer) GetIdxCount(dt archives.DocumentType) (int64, error) {
	return -1, errors.New("not yet implemented")
}

func (p *PGXArchiveIndexer) AddArtistRelease(artistDid int64, releaseDid int64, role IdxArtistRoleInRelease) error {
	sql := fmt.Sprintf("INSERT INTO %s (artist_did, release_did, role) VALUES ($1, $2, $3)", IdxArtistReleases)
	commandTag, err := p.dbp.Exec(context.Background(), sql, artistDid, releaseDid, role)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("artist release not added")
	}

	return nil
}

func (p *PGXArchiveIndexer) AddLabelRelease(labelDid int64, releaseDid int64) error {
	sql := fmt.Sprintf("INSERT INTO %s (label_did, release_did) VALUES ($1, $2)", IdxLabelReleases)
	commandTag, err := p.dbp.Exec(context.Background(), sql, labelDid, releaseDid)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("label release not added")
	}

	return nil
}

func (p *PGXArchiveIndexer) AddArtistSearchInfo(artistDid int64, as *discogstypes.ArtistSearchInfo) error {
	return errors.New("not yet implemented")
}

func (p *PGXArchiveIndexer) AddGenre(rg *discogstypes.ReleaseGenreEntry) error {
	return errors.New("not yet implemented")
}

func (p *PGXArchiveIndexer) AddStyle(rs *discogstypes.ReleaseStyleEntry) error {
	return errors.New("not yet implemented")
}

func (p *PGXArchiveIndexer) GetArtistReleases(artistDid int64) ([]IdxDiscogsReleaseByArtist, error) {
	sql := fmt.Sprintf("SELECT release_did, role FROM %s WHERE artist_did = $1", IdxArtistReleases)
	rows, err := p.dbp.Query(context.Background(), sql, artistDid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	artistReleases := make([]IdxDiscogsReleaseByArtist, 128)
	for rows.Next() {
		var releaseDid int64
		var role int32
		err = rows.Scan(&releaseDid, &role)
		if err != nil {
			return nil, err
		}
		artistRelease := IdxDiscogsReleaseByArtist{
			ReleaseDid: releaseDid,
			Role:       role,
		}
		artistReleases = append(artistReleases, artistRelease)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return artistReleases, nil
}

func (p *PGXArchiveIndexer) GetRangeOfDocumentIds(idxName IdxName, start int64, end int64) ([]int64, error) {
	sql := fmt.Sprintf("SELECT did FROM %s WHERE did BETWEEN $1 AND $2", idxName)

	rows, err := p.dbp.Query(context.Background(), sql, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	count := 0
	dids := make([]int64, 0, BlockSize)
	for rows.Next() {
		var did int64
		err = rows.Scan(&did)
		if err != nil {
			return nil, err
		}
		dids = append(dids, did)
		count++
	}

	if rows.Err() != nil {
		return nil, err
	}

	return dids, nil
}

func (p *PGXArchiveIndexer) Begin() (DiscogsArchiveIndexerTx, error) {
	tx, err := p.dbp.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	ptx := newPGXIndexerTx(tx)
	return ptx, nil
}

func (p *PGXIndexerTx) Rollback() error {
	err := p.tx.Rollback(context.Background())
	return err
}

func (p *PGXIndexerTx) Commit() error {
	err := p.tx.Commit(context.Background())
	return err
}
