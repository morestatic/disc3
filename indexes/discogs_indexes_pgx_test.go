package indexes_test

import (
	"fmt"
	"testing"

	"deepsolutionsvn.com/disc/archives"
	"deepsolutionsvn.com/disc/indexes"
)

const TestPGXDBURL = "postgres://localhost:5432/discogs"

func TestPGXShouldGetStartAndEndPosForRelease(t *testing.T) {
	t.Skip("postgresql not currently being supported")

	indexer, err := indexes.NewPGXArchiveIndexer(TestPGXDBURL)
	if err != nil {
		t.Fatal("unable to connect to db")
	}
	defer indexer.Close()

	startPos, endPos, err := indexer.GetContentIdx(archives.Releases, 5)
	if err != nil {
		t.Fatal("failed to read start and end pos", err)
	}

	if startPos != 26452 {
		t.Fatal("incorrect startPos", startPos)
	}

	if endPos != 31746 {
		t.Fatal("incorrect endPos", endPos)
	}
}

func TestPGXShouldGetArtistReleases(t *testing.T) {
	t.Skip("postgresql not currently being supported")

	indexer, err := indexes.NewPGXArchiveIndexer(TestPGXDBURL)
	if err != nil {
		t.Fatal("unable to connect to db")
	}
	defer indexer.Close()

	artistDid := int64(5)
	artistReleases, err := indexer.GetArtistReleases(artistDid)
	if err != nil {
		t.Fatal("failed to get artist releases", err)
	}

	numReleases := len(artistReleases)
	if numReleases != 599 {
		t.Fatal("unexpected number of releases for artist", numReleases)
	}

	fmt.Println(len(artistReleases))
}

func TestPGXShouldGatherArtistReleasesByRole(t *testing.T) {
	t.Skip("postgresql not currently being supported")

	indexer, err := indexes.NewPGXArchiveIndexer(TestPGXDBURL)
	if err != nil {
		t.Fatal("unable to connect to db")
	}
	defer indexer.Close()

	artistDid := int64(5)
	artistReleases, err := indexer.GetArtistReleases(artistDid)
	if err != nil {
		t.Fatal("failed to get artist releases", err)
	}

	artistReleasesByRole, err := indexes.GroupArtistReleasesByRole(artistReleases)
	if err != nil {
		t.Fatal("error grouping artist releases", err)
	}

	numPrimaryRole := len(artistReleasesByRole[indexes.PrimaryRole])
	if numPrimaryRole != 154 {
		t.Fatal("unexpected number of releases as primary artist", numPrimaryRole)
	}

	numExtraRole := len(artistReleasesByRole[indexes.ExtraRole])
	if numExtraRole != 108 {
		t.Fatal("unexpected number of releases as extra artist", numExtraRole)
	}

	numTrackRole := len(artistReleasesByRole[indexes.TrackRole])
	if numTrackRole != 272 {
		t.Fatal("unexpected number of releases as track artist", numTrackRole)
	}
}

func TestPGXShouldGetReleasesInRangeOfFirstBlock(t *testing.T) {
	t.Skip("postgresql not currently being supported")

	indexer, err := indexes.NewPGXArchiveIndexer(TestPGXDBURL)
	if err != nil {
		t.Fatal("unable to connect to db")
	}
	defer indexer.Close()

	count := int64(1)
	start, end := indexes.GetBlockRange(count)
	releasesInRange, err := indexer.GetRangeOfDocumentIds(indexes.IdxReleases, start, end)
	if err != nil {
		t.Fatal("error getting range of document ids", err)
	}
	if len(releasesInRange) != 10000 {
		t.Fatal("unexpected count of release document ids", len(releasesInRange))
	}
	if releasesInRange[0] != 1 {
		t.Fatal("unexpected first release", releasesInRange[0])
	}
	if releasesInRange[19] != 20 {
		t.Fatal("unexpected first release", releasesInRange[19])
	}
}
