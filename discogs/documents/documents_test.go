package documents_test

import (
	"fmt"
	"testing"

	archives "deepsolutionsvn.com/disc/discogs/archives"
	commands "deepsolutionsvn.com/disc/discogs/commands"
	documents "deepsolutionsvn.com/disc/discogs/documents"
	indexes "deepsolutionsvn.com/disc/discogs/indexes"
)

const ExistingTestSQLiteDBURL = "file:../testdata/db/test_002.db"

func TestShouldConvertXmlToJson(t *testing.T) {
	f, err := archives.NewDiscogsFileReader(archives.Releases, "002", "../testdata/drops")
	if err != nil {
		t.Fatal("unable to read archive file")
	}

	startPos := int64(877)
	endPos := int64(2174)

	xml, _, err := archives.GetContent(f, startPos, endPos)
	if err != nil {
		t.Fatal("unable to get xml content", err)
	}

	json, err := documents.GetJson(xml)
	if err != nil {
		t.Fatal("error getting json", err)
	}

	if len(json) != 892 {
		t.Fatal("incorrect length for json", len(json))
	}
}

func TestShouldCalcBuckets(t *testing.T) {
	l1, l2, l3 := documents.CalcBuckets(20030201)

	if l1 != 20 {
		t.Fatal("incorrect l1 bucket", l1)
	}

	if l2 != 3 {
		t.Fatal("incorrect l2 bucket", l2)
	}

	if l3 != 2 {
		t.Fatal("incorrect l3 bucket", l3)
	}
}

// func TestShouldCalcS2Buckets(t *testing.T) {
// 	l1, l2 := documents.CalcS2Buckets(19131221)

// 	if l1 != 19 {
// 		t.Fatal("incorrect l1 bucket", l1)
// 	}

// 	if l2 != 131 {
// 		t.Fatal("incorrect l2 bucket", l2)
// 	}
// }

func TestShouldMakeFolderForArtist(t *testing.T) {
	artistDid := int64(5)

	artistFolderName, err := documents.MakeArtistFolderName("../testdata/staging", artistDid)
	if err != nil {
		t.Fatal("error making artist folder name", artistFolderName)
	}

	if artistFolderName != "../testdata/staging/A/00/00/00/00000005" {
		t.Fatal("incorrect artist folder name", artistFolderName)
	}
}

// func TestShouldMakeS2FolderForArtist(t *testing.T) {
// 	artistDid := int64(5)

// 	artistFolderName, err := documents.MakeArtistS2FolderName("../testdata/staging", artistDid)
// 	if err != nil {
// 		t.Fatal("error making artist folder name", artistFolderName)
// 	}

// 	if artistFolderName != "../testdata/staging/A/000/000/00000005" {
// 		t.Fatal("incorrect artist folder name", artistFolderName)
// 	}
// }

func getTestJsonForArtist(t *testing.T) (string, error) {
	f, err := archives.NewDiscogsFileReader(archives.Artists, "002", "../testdata/drops")
	if err != nil {
		t.Fatal("unable to read archive file")
	}

	startPos := int64(2897)
	endPos := int64(3672)

	xml, _, err := archives.GetContent(f, startPos, endPos)
	if err != nil {
		t.Fatal("unable to get xml content", err)
	}

	artistJson, err := documents.GetJson(xml)
	if err != nil {
		t.Fatal("error getting json", err)
	}

	if len(artistJson) != 683 {
		t.Fatal("incorrect length for json", len(artistJson))
	}

	return artistJson, nil
}

func TestShouldWriteArtistInfoFile(t *testing.T) {
	artistDid := int64(50)

	artistJson, err := getTestJsonForArtist(t)
	if err != nil {
		t.Fatal("problem reading test json", err)
	}

	filename, err := documents.WriteArtistInfoToStaging("../testdata/staging", artistDid, artistJson)
	if err != nil {
		t.Fatal("failed to write artist info to staging", err)
	}

	if filename != "../testdata/staging/A/00/00/00/00000050/artist_info.json" {
		t.Fatalf("unexpected artist_info filename: %s\n", filename)
	}
}

func TestShouldWriteAndReadArtistReleasesJsonFromStaging(t *testing.T) {
	artistDid := int64(50)

	d, err := archives.NewDiscogsFileDecoder(archives.Releases, "002", "../testdata/drops")
	if err != nil || d == nil {
		t.Fatal(fmt.Errorf("failed to create new decoder: %w", err))
	}
	defer d.Close()

	i, err := indexes.NewSQLiteDiscogsArchiveIndexer(ExistingTestSQLiteDBURL)
	if err != nil {
		t.Fatal(fmt.Errorf("unable to connect to indexer: %w", err))
	}

	err = commands.SaveArtistReleasesAsJsonDoc(artistDid, "../testdata/staging", i)
	if err != nil {
		t.Fatal(fmt.Errorf("failed to save artist releases json: %w", err))
	}

	artistReleasesByRole, err := documents.ReadArtistReleasesJsonFromStaging("../testdata/staging", artistDid)
	if err != nil {
		t.Fatal("failed to read artist releases from json library staging", err)
	}

	numPrimaryReleases := len(artistReleasesByRole[indexes.PrimaryRole])
	if numPrimaryReleases != 90 {
		t.Fatal("unexpected number of releases as primary artist", numPrimaryReleases)
	}

	numExtraReleases := len(artistReleasesByRole[indexes.ExtraRole])
	if numExtraReleases != 7 {
		t.Fatal("unexpected number of releases as extra artist", numExtraReleases)
	}

	numTrackReleases := len(artistReleasesByRole[indexes.TrackRole])
	if numTrackReleases != 70 {
		t.Fatal("unexpected number of releases as track artist", numTrackReleases)
	}
}
