package indexes

import (
	"errors"
	"os"

	discogstypes "deepsolutionsvn.com/disc/discogs/types"

	archives "deepsolutionsvn.com/disc/discogs/archives"
)

type IdxName string

const (
	Unknown              = iota
	IdxArtists           = "artists"
	IdxLabels            = "labels"
	IdxReleases          = "releases"
	IdxMasters           = "masters"
	IdxArtistReleases    = "artist_releases"
	IdxLabelReleases     = "label_releases"
	IdxArtistsSearchInfo = "artists_search_info"
	IdxReleaseGenres     = "release_genres"
	IdxReleaseStyles     = "release_styles"
)

type IdxArtistRoleInRelease = int32

const (
	IdxPrimaryRole = 1
	IdxExtraRole   = 2
	IdxTrackRole   = 3
)

type ArtistRoleInRelease = string

const (
	PrimaryRole = "primary"
	ExtraRole   = "extra"
	TrackRole   = "track"
)

type DiscogsArchiveIndexer interface {
	GetContentIdx(dt archives.DocumentType, did int64) (int64, int64, error)
	AddWithContentIdx(dt archives.DocumentType, did int64, startPos int64, endPos int64) error
	UpdateContentIdx(dt archives.DocumentType, did int64, startPos int64, endPos int64) error
	AddIdxCount(dt archives.DocumentType, count int64) error
	GetIdxCount(dt archives.DocumentType) (int64, error)
	AddArtistRelease(artistDid int64, releaseDid int64, role IdxArtistRoleInRelease) error
	AddLabelRelease(labelDid int64, releaseDid int64) error
	GetArtistReleases(artistDid int64) ([]IdxDiscogsReleaseByArtist, error)
	GetRangeOfDocumentIds(idxName IdxName, start int64, end int64) ([]int64, error)
	AddArtistSearchInfo(artistDid int64, as *discogstypes.ArtistSearchInfo) error
	AddGenre(rg *discogstypes.ReleaseGenreEntry) error
	AddStyle(rs *discogstypes.ReleaseStyleEntry) error
	Close()
	Begin() (DiscogsArchiveIndexerTx, error)
}

type DiscogsArchiveIndexerTx interface {
	Rollback() error
	Commit() error
}

type IdxDiscogsReleaseByArtist struct {
	ReleaseDid int64
	Role       IdxArtistRoleInRelease
}

type DiscogsReleasesByArtistByRole map[ArtistRoleInRelease][]int64
type DiscogsReleaseDidSet map[int64]bool

var IdxUrl = "postgres://localhost:5432/discogs"

func (dbName IdxName) String() string {
	switch dbName {
	case IdxArtists:
		return "artists"
	case IdxLabels:
		return "labels"
	case IdxReleases:
		return "releases"
	case IdxMasters:
		return "masters"
	case IdxArtistReleases:
		return "artist_releases"
	case IdxLabelReleases:
		return "label_releases"
	default:
		return "unknown"
	}
}

func GetDefaultConnUrl() string {
	return os.Getenv("DATABASE_URL")
}

func GroupArtistReleasesByRole(artistReleases []IdxDiscogsReleaseByArtist) (DiscogsReleasesByArtistByRole, error) {
	byRole := make(DiscogsReleasesByArtistByRole, 3)
	hasDid := make(map[string]DiscogsReleaseDidSet)
	hasDid[PrimaryRole] = make(DiscogsReleaseDidSet, 1024*64)
	hasDid[ExtraRole] = make(DiscogsReleaseDidSet, 1024*64)
	hasDid[TrackRole] = make(DiscogsReleaseDidSet, 1024*64)
	for _, artistRelease := range artistReleases {
		did := artistRelease.ReleaseDid
		switch artistRelease.Role {
		case IdxPrimaryRole:
			if !hasDid[PrimaryRole][did] {
				byRole[PrimaryRole] = append(byRole[PrimaryRole], did)
				hasDid[PrimaryRole][did] = true
			}
		case IdxExtraRole:
			if !hasDid[ExtraRole][did] {
				byRole[ExtraRole] = append(byRole[ExtraRole], did)
				hasDid[ExtraRole][did] = true
			}
		case IdxTrackRole:
			if !hasDid[TrackRole][did] {
				byRole[TrackRole] = append(byRole[TrackRole], did)
				hasDid[TrackRole][did] = true
			}
		default:
			return nil, errors.New("unknown role for artist")
		}
	}
	return byRole, nil
}
