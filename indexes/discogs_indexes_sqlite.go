package indexes

import (
	"context"
	"errors"
	"fmt"

	// "crawshaw.io/sqlite"
	// "crawshaw.io/sqlite/sqlitex"
	// "zombiezen.com/go/sqlite"
	// "zombiezen.com/go/sqlite/sqlitex"

	"github.com/dislabsvn/sqlite"
	"github.com/dislabsvn/sqlite/sqlitex"

	discogstypes "deepsolutionsvn.com/disc/types/discogs"

	"deepsolutionsvn.com/disc/archives"
)

type SQLiteDiscogsArchiveIndexer struct {
	dbp      *sqlitex.Pool
	activeTx *SQLiteDiscogsIndexerTx
}

type SQLiteDiscogsIndexerTx struct {
	conn *sqlite.Conn
	p    *SQLiteDiscogsArchiveIndexer
}

func NewSQLiteDiscogsArchiveIndexer(connString string) (*SQLiteDiscogsArchiveIndexer, error) {
	if connString == "" {
		connString = GetDefaultConnUrl()
	}

	fmt.Println(connString)
	initScript := "PRAGMA journal_mode = MEMORY; PRAGMA synchronous = 0;"
	// initScript := "PRAGMA synchronous = 1;"
	// initScript := "PRAGMA synchronous = 1;"
	flags := sqlite.SQLITE_OPEN_EXCLUSIVE | sqlite.SQLITE_OPEN_URI | sqlite.SQLITE_OPEN_READWRITE
	dbpool, err := sqlitex.OpenInitWithOpts(context.Background(), connString, flags, 16, initScript, sqlitex.ExecScriptOpts{SkipTx: true})
	// dbpool, err := sqlitex.OpenInit(context.Background(), connString, flags, 16, initScript)
	// dbpool, err := sqlitex.Open(connString, 0, 16)
	if err != nil {
		return nil, err
	}

	connector := &SQLiteDiscogsArchiveIndexer{
		dbp: dbpool,
	}

	return connector, nil
}

func newSQLiteDiscogsIndexerTx(p *SQLiteDiscogsArchiveIndexer) *SQLiteDiscogsIndexerTx {
	ptx := &SQLiteDiscogsIndexerTx{
		p: p,
	}
	return ptx
}

func (p *SQLiteDiscogsArchiveIndexer) Close() {
	p.dbp.Close()
}

func (p *SQLiteDiscogsArchiveIndexer) GetContentIdx(dt archives.DocumentType, did int64) (int64, int64, error) {

	conn := p.dbp.Get(context.Background())
	if conn == nil {
		return -1, -1, errors.New("unable to get connection from pool")
	}
	defer p.dbp.Put(conn)

	sql := fmt.Sprintf("SELECT start_pos, end_pos FROM %s WHERE did = $did", dt)
	stmt, err := conn.Prepare(sql)
	if err != nil {
		return -1, -1, err
	}
	defer stmt.Reset()

	stmt.SetInt64("$did", did)

	hasRow, err := stmt.Step()
	if err != nil {
		return -1, -1, err
	}
	if !hasRow {
		return -1, -1, fmt.Errorf("entity with did: %d not found", did)
	}

	startPos := stmt.GetInt64("start_pos")
	endPos := stmt.GetInt64("end_pos")

	return startPos, endPos, nil
}

func (p *SQLiteDiscogsArchiveIndexer) AddWithContentIdx(dt archives.DocumentType, did int64, startPos int64, endPos int64) error {

	conn := p.dbp.Get(context.Background())
	if conn == nil {
		return errors.New("unable to get connection from pool")
	}
	defer p.dbp.Put(conn)
	sql := fmt.Sprintf("INSERT INTO %s (did, start_pos, end_pos) VALUES ($did, $start_pos, $end_pos)", dt)
	// fmt.Println(sql)
	stmt, err := conn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Reset()

	stmt.SetInt64("$did", did)
	stmt.SetInt64("$start_pos", startPos)
	stmt.SetInt64("$end_pos", endPos)

	_, err = stmt.Step()
	if err != nil {
		return err
	}

	// if !hasRow {
	// 	return fmt.Errorf("artist with did: %d not found", did)
	// }

	return nil
}

func (p *SQLiteDiscogsArchiveIndexer) UpdateContentIdx(dt archives.DocumentType, did int64, startPos int64, endPos int64) error {
	// sql := fmt.Sprintf("UPDATE %s SET start_pos = $2, end_pos = $3 WHERE did = $1", dt)

	// commandTag, err := p.dbp.Exec(context.Background(), sql, did, startPos, endPos)
	// if err != nil {
	// 	return err
	// }
	// if commandTag.RowsAffected() != 1 {
	// 	return errors.New("content idx not updated")
	// }
	return nil
}

func (p *SQLiteDiscogsArchiveIndexer) AddIdxCount(dt archives.DocumentType, count int64) error {

	conn := p.dbp.Get(context.Background())
	if conn == nil {
		return errors.New("unable to get connection from pool")
	}
	defer p.dbp.Put(conn)

	propertyName := fmt.Sprintf("%s_total", dt.String())
	sql := "INSERT INTO stats (property, value) VALUES ($property, $value)"
	stmt, err := conn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Reset()

	stmt.SetText("$property", propertyName)
	stmt.SetInt64("$value", count)

	_, err = stmt.Step()
	if err != nil {
		return err
	}

	// if !hasRow {
	// 	return fmt.Errorf("artist with did: %d not found", did)
	// }

	return nil
}

func (p *SQLiteDiscogsArchiveIndexer) GetIdxCount(dt archives.DocumentType) (int64, error) {
	conn := p.dbp.Get(context.Background())
	if conn == nil {
		return -1, errors.New("unable to get connection from pool")
	}
	defer p.dbp.Put(conn)

	propertyName := fmt.Sprintf("%s_total", dt.String())

	sql := "SELECT property, value FROM stats WHERE property = $property"
	stmt, err := conn.Prepare(sql)
	if err != nil {
		return -1, err
	}
	defer stmt.Reset()

	stmt.SetText("$property", propertyName)

	hasRow, err := stmt.Step()
	if err != nil {
		return -1, err
	}
	if !hasRow {
		return -1, fmt.Errorf("property with %s not found", propertyName)
	}

	total := stmt.GetInt64("value")

	return total, nil
}

func (p *SQLiteDiscogsArchiveIndexer) AddArtistRelease(artistDid int64, releaseDid int64, role IdxArtistRoleInRelease) error {

	var conn *sqlite.Conn
	if p.activeTx != nil {
		conn = p.activeTx.conn
	} else {
		conn = p.dbp.Get(context.Background())
		if conn == nil {
			return errors.New("unable to get connection from pool")
		}
		defer p.dbp.Put(conn)
	}

	sql := fmt.Sprintf("INSERT INTO %s (artist_did, release_did, role) VALUES ($artist_did, $release_did, $role)", IdxArtistReleases)
	stmt, err := conn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.ClearBindings()
	defer stmt.Reset()

	stmt.SetInt64("$artist_did", artistDid)
	stmt.SetInt64("$release_did", releaseDid)
	stmt.SetInt64("$role", int64(role))

	_, err = stmt.Step()
	if err != nil {
		return err
	}

	return nil
}

func (p *SQLiteDiscogsArchiveIndexer) AddLabelRelease(labelDid int64, releaseDid int64) error {
	// sql := fmt.Sprintf("INSERT INTO %s (label_did, release_did) VALUES ($1, $2)", IdxLabelReleases)
	// commandTag, err := p.dbp.Exec(context.Background(), sql, labelDid, releaseDid)
	// if err != nil {
	// 	return err
	// }
	// if commandTag.RowsAffected() != 1 {
	// 	return errors.New("label release not added")
	// }
	return nil
}

func (p *SQLiteDiscogsArchiveIndexer) AddArtistSearchInfo(artistDid int64, as *discogstypes.ArtistSearchInfo) error {
	conn := p.dbp.Get(context.Background())
	if conn == nil {
		return errors.New("unable to get connection from pool")
	}
	defer p.dbp.Put(conn)
	sql := fmt.Sprintf("INSERT INTO %s (artist_did, name, realname, is_group) VALUES ($artist_did, $name, $realname, $is_group)", IdxArtistsSearchInfo)
	stmt, err := conn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Reset()

	stmt.SetInt64("$artist_did", artistDid)
	stmt.SetText("$name", as.Name)
	stmt.SetText("$realname", as.RealName)

	var isGroup int64 = 0
	if as.IsGroup {
		isGroup = 1
	}
	stmt.SetInt64("$is_group", isGroup)

	_, err = stmt.Step()
	if err != nil {
		return err
	}

	return nil
}

func (p *SQLiteDiscogsArchiveIndexer) AddGenre(rg *discogstypes.ReleaseGenreEntry) error {
	conn := p.dbp.Get(context.Background())
	if conn == nil {
		return errors.New("unable to get connection from pool")
	}
	defer p.dbp.Put(conn)
	sql := fmt.Sprintf("INSERT INTO %s (genre) VALUES ($genre) ON CONFLICT DO NOTHING", IdxReleaseGenres)
	stmt, err := conn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Reset()

	stmt.SetText("$genre", rg.Genre)

	_, err = stmt.Step()
	if err != nil {
		return err
	}

	return nil
}

func (p *SQLiteDiscogsArchiveIndexer) AddStyle(rs *discogstypes.ReleaseStyleEntry) error {
	conn := p.dbp.Get(context.Background())
	if conn == nil {
		return errors.New("unable to get connection from pool")
	}
	defer p.dbp.Put(conn)
	sql := fmt.Sprintf("INSERT INTO %s (style) VALUES ($style) ON CONFLICT DO NOTHING", IdxReleaseStyles)
	stmt, err := conn.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Reset()

	stmt.SetText("$style", rs.Style)

	_, err = stmt.Step()
	if err != nil {
		return err
	}

	return nil
}

func (p *SQLiteDiscogsArchiveIndexer) GetArtistReleases(artistDid int64) ([]IdxDiscogsReleaseByArtist, error) {

	conn := p.dbp.Get(context.Background())
	if conn == nil {
		return nil, errors.New("unable to get connection from pool")
	}
	defer p.dbp.Put(conn)

	sql := fmt.Sprintf("SELECT release_did, role FROM %s WHERE artist_did = $artist_did", IdxArtistReleases)
	stmt, err := conn.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Reset()

	stmt.SetInt64("$artist_did", artistDid)

	artistReleases := make([]IdxDiscogsReleaseByArtist, 0)
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return nil, err
		}
		if !hasRow {
			break
		}

		releaseDid := stmt.GetInt64("release_did")
		role := stmt.GetInt64("role")

		artistRelease := IdxDiscogsReleaseByArtist{
			ReleaseDid: releaseDid,
			Role:       int32(role),
		}
		artistReleases = append(artistReleases, artistRelease)
	}

	return artistReleases, nil
}

func (p *SQLiteDiscogsArchiveIndexer) GetRangeOfDocumentIds(idxName IdxName, start int64, end int64) ([]int64, error) {

	conn := p.dbp.Get(context.Background())
	if conn == nil {
		return nil, errors.New("unable to get connection from pool")
	}
	defer p.dbp.Put(conn)

	sql := fmt.Sprintf("SELECT did FROM %s WHERE did BETWEEN $start_did AND $end_did", idxName)
	stmt, err := conn.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Reset()

	stmt.SetInt64("$start_did", start)
	stmt.SetInt64("$end_did", end)

	dids := make([]int64, 0, BlockSize)
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return nil, err
		}
		if !hasRow {
			break
		}

		did := stmt.GetInt64("did")

		dids = append(dids, did)
	}

	return dids, nil
}

func (p *SQLiteDiscogsArchiveIndexer) Begin() (DiscogsArchiveIndexerTx, error) {
	ptx := newSQLiteDiscogsIndexerTx(p)

	conn := ptx.p.dbp.Get(context.Background())
	if conn == nil {
		return nil, errors.New("unable to get connection from pool")
	}

	ptx.conn = conn
	p.activeTx = ptx

	err := sqlitex.ExecTransient(ptx.conn, "BEGIN TRANSACTION;", nil)
	if err != nil {
		ptx.p.dbp.Put(ptx.conn)
		ptx.p.activeTx = nil
		return nil, fmt.Errorf("unable to begin transaction: %w", err)
	}

	return ptx, nil
}

func (ptx *SQLiteDiscogsIndexerTx) Rollback() error {

	var err error

	// don't try to rollback if no activeTx
	if ptx.p != nil && ptx.p.activeTx != nil {
		defer func() {
			ptx.p.dbp.Put(ptx.conn)
			ptx.p.activeTx = nil
		}()

		err = sqlitex.ExecTransient(ptx.conn, "ROLLBACK TRANSACTION;", nil)
		if err != nil {
			return fmt.Errorf("unable to rollback transaction: %w", err)
		}
	}

	return nil
}

func (ptx *SQLiteDiscogsIndexerTx) Commit() error {

	defer func() {
		ptx.p.dbp.Put(ptx.conn)
		ptx.p.activeTx = nil
	}()

	err := sqlitex.ExecTransient(ptx.conn, "COMMIT TRANSACTION;", nil)
	if err != nil {
		return fmt.Errorf("unable to commit transaction: %w", err)
	}

	return nil
}
