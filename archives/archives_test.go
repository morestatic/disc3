package archives_test

import (
	"fmt"
	"os"
	"testing"

	"deepsolutionsvn.com/disc/archives"
	discogstypes "deepsolutionsvn.com/disc/types/discogs"
)

func TestShouldMakeReleasesArchiveName(t *testing.T) {
	result := archives.MakeDefaultArchiveName(archives.Releases, "20211001", ".")
	expected := "./discogs_20211001_releases.xml"
	if result != expected {
		t.Fatal("mismatch when making archive name", result, expected)
	}
}

func TestShouldMakeArtistsArchiveName(t *testing.T) {
	result := archives.MakeDefaultArchiveName(archives.Artists, "20211001", ".")
	expected := "./discogs_20211001_artists.xml"
	if result != expected {
		t.Fatal("mismatch when making archive name", result, expected)
	}
}

func TestShouldMakeLabelsArchiveName(t *testing.T) {
	result := archives.MakeDefaultArchiveName(archives.Labels, "20211001", "./discogs1")
	expected := "./discogs1/discogs_20211001_labels.xml"
	if result != expected {
		t.Fatal("mismatch when making archive name", result, expected)
	}
}

func TestShouldMakeMastersArchiveName(t *testing.T) {
	result := archives.MakeDefaultArchiveName(archives.Masters, "20211001", ".")
	expected := "./discogs_20211001_masters.xml"
	if result != expected {
		t.Fatal("mismatch when making archive name", result, expected)
	}
}

func TestShouldCreateNewFileDecoder(t *testing.T) {
	d, err := archives.NewDiscogsFileDecoder(archives.Releases, "002", "../testdata/drops")
	if err != nil || d == nil {
		t.Fatal("failed to create new decoder", err)
	}
	defer d.Close()

	if d.GetFile() == nil {
		t.Fatal("decoder file not open")
	}
}

func TestShouldFailToCreateNewFileDecoder(t *testing.T) {
	decoder, err := archives.NewDiscogsFileDecoder(archives.Releases, "2021100", "../discogs")
	if err == nil || decoder != nil {
		t.Fatal("created decoder when shouldn't have", err)
	}
}

func TestShouldFindFirstElementInReleasesScan(t *testing.T) {
	d, err := archives.NewDiscogsFileDecoder(archives.Releases, "001", "../testdata/drops")
	if err != nil || d == nil {
		t.Fatal("failed to create new decoder")
	}
	defer d.Close()

	var rx *discogstypes.MinimalReleaseXml = nil
	var progress archives.DiscogsScanProgress

	_, count, err := d.Scan(archives.Releases, "release", discogstypes.MinimalReleaseXml{}, func(d *archives.DiscogsFileDecoder, e interface{}, o interface{}, p archives.DiscogsScanProgress) (bool, error) {
		rx = o.(*discogstypes.MinimalReleaseXml)
		progress = p
		return true, nil
	})

	if err != nil {
		t.Fatal("scan of first release failed", err)
	}

	if count != 1 {
		t.Fatal("unexpected count", count)
	}

	if progress.StartPos != 10 {
		t.Fatal("unexpected start pos", progress.StartPos)
	}

	if rx.Id != 20207 {
		t.Fatal("first release id not equal 606975", rx.Id)
	}

}

func TestShouldFindFirstElementInArtistsScan(t *testing.T) {
	d, err := archives.NewDiscogsFileDecoder(archives.Artists, "002", "../testdata/drops")
	if err != nil || d == nil {
		t.Fatal("failed to create new decoder")
	}
	defer d.Close()

	var ax *discogstypes.MinimalArtistXml = nil
	var progress archives.DiscogsScanProgress

	_, count, err := d.Scan(archives.Artists, "artist", discogstypes.MinimalArtistXml{}, func(d *archives.DiscogsFileDecoder, e interface{}, o interface{}, p archives.DiscogsScanProgress) (bool, error) {
		ax = o.(*discogstypes.MinimalArtistXml)
		progress = p
		return true, nil
	})

	if err != nil {
		t.Fatal("scan of first artist failed")
	}

	if count != 1 {
		t.Fatal("unexpected count", count)
	}

	if progress.StartPos != 9 {
		t.Fatal("unexpected start pos")
	}

	if ax.Id != 20 {
		t.Fatal("first artist id not equal 1", ax.Id)
	}
}

func TestShouldFindSecondElementInArtistsScan(t *testing.T) {
	d, err := archives.NewDiscogsFileDecoder(archives.Artists, "002", "../testdata/drops")
	if err != nil || d == nil {
		t.Fatal("failed to create new decoder")
	}
	defer d.Close()

	var ax *discogstypes.MinimalArtistXml = nil
	var progress archives.DiscogsScanProgress

	_, count, err := d.Scan(archives.Artists, "artist", discogstypes.MinimalArtistXml{}, func(d *archives.DiscogsFileDecoder, e interface{}, o interface{}, p archives.DiscogsScanProgress) (bool, error) {
		ax = o.(*discogstypes.MinimalArtistXml)
		progress = p
		if p.Count == 2 {
			return true, nil
		} else {
			return false, nil
		}
	})

	if err != nil {
		t.Fatal("scan of first artist failed", err)
	}

	if count != 2 {
		t.Fatal("unexpected count", count)
	}

	if progress.StartPos != 877 {
		t.Fatal("unexpected start pos", progress.StartPos)
	}

	if ax.Id != 30 {
		t.Fatal("second artist id not equal 30", ax.Id)
	}
}

func TestShouldReadReleaseXmlFromArchive(t *testing.T) {
	d, err := archives.NewDiscogsFileDecoder(archives.Releases, "001", "../testdata/drops")
	if err != nil || d == nil {
		t.Fatal("failed to create new decoder", err)
	}
	defer d.Close()

	var rx *discogstypes.MinimalReleaseXml = nil
	var progress archives.DiscogsScanProgress

	_, count, err := d.Scan(archives.Releases, "release", discogstypes.MinimalReleaseXml{}, func(d *archives.DiscogsFileDecoder, e interface{}, o interface{}, p archives.DiscogsScanProgress) (bool, error) {
		rx = o.(*discogstypes.MinimalReleaseXml)
		progress = p

		if p.Count == 5 {
			fmt.Println(rx)
			fmt.Println(p)
			return true, nil
		} else {
			return false, nil
		}
	})

	if err != nil {
		t.Fatal("failed scanning releases", err)
	}

	if count != 5 {
		t.Fatal("unexpected count", count)
	}

	if progress.StartPos != 16974 {
		t.Fatal("unexpected start pos", progress.StartPos)
	}

	if rx.Id != 11334732 {
		t.Fatal("unexpected release id", rx.Id)
	}

	_, len, err := archives.GetContent(d.GetFile(), progress.StartPos, progress.EndPos)
	if err != nil {
		t.Fatal("error reading xml content", err)
	}

	if len != 3085 {
		t.Fatal("unexpected xml size", len)
	}
}

func TestShouldCreateReleasesXmlWithNoReleases(t *testing.T) {
	filename := "../testdata/test_releases_empty.xml"
	w, err := archives.NewDiscogsFileEncoder(archives.Releases, filename)
	if err != nil {
		t.Fatal("failed to create new encoder", err)
	}
	w.Close()
	contents, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal("failed to read test file contents", err)
	}
	if len(contents) != 22 {
		t.Fatal("unexpected length of test file contents", len(contents))
	}
}
