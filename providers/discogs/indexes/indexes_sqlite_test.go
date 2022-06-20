package indexes_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	archives "deepsolutionsvn.com/disc/providers/discogs/archives"
	indexes "deepsolutionsvn.com/disc/providers/discogs/indexes"
)

const EmptyTestSQLiteDB = "../testdata/db/test_empty.db"
const ExistingTestSQLiteDB = "../testdata/db/test_002.db"
const NewTestSQLiteDB = "../testdata/db/sqlite/test_new.db"
const TestSQLiteDB = "../testdata/db/test_existing.db"

const NewTestSQLiteDBURL = "file:../testdata/db/sqlite/test_new.db"
const ExistingTestSQLiteDBURL = "file:../testdata/db/test_002.db"
const EmptyTestSQLiteDBURL = "file:../testdata/db/test_empty.db"
const TestSQLiteDBURL = "file:../testdata/db/test_existing.db"

func testDBExists(testDB string) (bool, error) {
	if _, err := os.Stat(testDB); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func removeTestDB(testDB string) error {
	rmCmd := exec.Command("rm", "-f", testDB)
	err := rmCmd.Run()
	return err
}

func makeTestDB(sourceDB string, destDB string) error {
	exists, err := testDBExists(destDB)
	if err != nil {
		return fmt.Errorf("unable to check whether the test db exists: %s, %w", destDB, err)
	}

	if exists {
		err = removeTestDB(destDB)
		if err != nil {
			return fmt.Errorf("unable to remove existing test db: %s, %w", destDB, err)
		}
	}

	copyCmd := exec.Command("cp", "-f", sourceDB, destDB)
	err = copyCmd.Run()
	if err != nil {
		return fmt.Errorf("unable to run copy cmd: %w", err)
	}

	return nil
}

func TestSQLiteShouldAddAndGetContentIdx(t *testing.T) {
	err := makeTestDB(EmptyTestSQLiteDB, NewTestSQLiteDB)
	if err != nil {
		t.Fatal(err)
	}

	indexer, err := indexes.NewSQLiteDiscogsArchiveIndexer(NewTestSQLiteDB)
	if err != nil {
		t.Fatal("unable to connect to db")
	}
	defer indexer.Close()

	did := int64(2784)
	startPos := int64(10)
	endPos := int64(2548)

	err = indexer.AddWithContentIdx(archives.Releases, did, startPos, endPos)
	if err != nil {
		t.Fatal("failed to add release content idx", err)
	}

	startPos, endPos, err = indexer.GetContentIdx(archives.Releases, 2784)
	if err != nil {
		t.Fatal("failed to read start and end pos", err)
	}

	if startPos != 10 {
		t.Fatal("incorrect startPos", startPos)
	}

	if endPos != 2548 {
		t.Fatal("incorrect endPos", endPos)
	}
}

func TestSQLiteShouldAddAndGetArtistRelease(t *testing.T) {
	err := makeTestDB(EmptyTestSQLiteDB, NewTestSQLiteDB)
	if err != nil {
		t.Fatal(err)
	}

	indexer, err := indexes.NewSQLiteDiscogsArchiveIndexer(NewTestSQLiteDBURL)
	if err != nil {
		t.Fatal("unable to connect to db")
	}
	defer indexer.Close()

	artistDid := int64(1)
	releaseDid := int64(100)
	role := indexes.IdxPrimaryRole

	err = indexer.AddArtistRelease(artistDid, releaseDid, int32(role))
	if err != nil {
		t.Fatal("failed to add artist release", err)
	}

	artistReleases, err := indexer.GetArtistReleases(artistDid)
	if err != nil {
		t.Fatal("failed to get artist releases", err)
	}

	if len(artistReleases) != 1 {
		t.Fatal("unexpected number of artist releases")
	}

	if artistReleases[0].ReleaseDid != 100 {
		t.Fatal("unexpected release for artist")
	}

	if artistReleases[0].Role != indexes.IdxPrimaryRole {
		t.Fatal("unexpected release role for artist")
	}
}

func TestSQLiteShouldGatherArtistReleasesByRole(t *testing.T) {
	err := makeTestDB(ExistingTestSQLiteDB, TestSQLiteDB)
	if err != nil {
		t.Fatal(err)
	}

	indexer, err := indexes.NewSQLiteDiscogsArchiveIndexer(TestSQLiteDBURL)
	if err != nil {
		t.Fatal("unable to connect to db")
	}
	defer indexer.Close()

	artistDid := int64(50)
	artistReleases, err := indexer.GetArtistReleases(artistDid)
	if err != nil {
		t.Fatal("failed to get artist releases", err)
	}

	artistReleasesByRole, err := indexes.GroupArtistReleasesByRole(artistReleases)
	if err != nil {
		t.Fatal("error grouping artist releases", err)
	}

	numPrimaryRole := len(artistReleasesByRole[indexes.PrimaryRole])
	if numPrimaryRole != 90 {
		t.Fatal("unexpected number of releases as primary artist", numPrimaryRole)
	}

	numExtraRole := len(artistReleasesByRole[indexes.ExtraRole])
	if numExtraRole != 7 {
		t.Fatal("unexpected number of releases as extra artist", numExtraRole)
	}

	numTrackRole := len(artistReleasesByRole[indexes.TrackRole])
	if numTrackRole != 70 {
		t.Fatal("unexpected number of releases as track artist", numTrackRole)
	}
}

func TestSQLiteShouldGetReleasesInRangeOfFirstBlock(t *testing.T) {
	err := makeTestDB(ExistingTestSQLiteDB, TestSQLiteDB)
	if err != nil {
		t.Fatal(err)
	}

	indexer, err := indexes.NewSQLiteDiscogsArchiveIndexer(TestSQLiteDBURL)
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
	if len(releasesInRange) != 84 {
		t.Fatal("unexpected count of release document ids", len(releasesInRange))
	}
	if releasesInRange[0] != 3 {
		t.Fatal("unexpected first release", releasesInRange[0])
	}
	if releasesInRange[19] != 135 {
		t.Fatal("unexpected first release", releasesInRange[19])
	}
}
