package archives_test

import (
	"testing"

	"deepsolutionsvn.com/disc/providers/musicbrainz/archives"
)

func TestShouldGetMusicBrainzArchivePathForDocType(t *testing.T) {
	cases := []struct {
		testName     string
		docType      archives.DocumentType
		dropBasePath string
		dropVersion  string
		expectedPath string
	}{
		{
			testName:     "Artists",
			docType:      archives.Artists,
			dropBasePath: "./MusicData/mb",
			dropVersion:  "20220501",
			expectedPath: "./MusicData/mb/20220501/A/mbdump/artist.json",
		},
		{
			testName:     "Releases",
			docType:      archives.Releases,
			dropBasePath: "./MusicData/mb",
			dropVersion:  "20220501",
			expectedPath: "./MusicData/mb/20220501/R/mbdump/release.json",
		},
		{
			testName:     "Labels",
			docType:      archives.Labels,
			dropBasePath: "./MusicData/mb",
			dropVersion:  "20220501",
			expectedPath: "./MusicData/mb/20220501/L/mbdump/label.json",
		},
		{
			testName:     "Masters",
			docType:      archives.Masters,
			dropBasePath: "./MusicData/mb",
			dropVersion:  "20220501",
			expectedPath: "./MusicData/mb/20220501/M/mbdump/release_group.json",
		},
	}

	for _, tc := range cases {
		t.Run(tc.testName, func(t *testing.T) {
			path := archives.MakeDefaultMusicBrainzArchiveName(tc.docType, tc.dropBasePath, tc.dropVersion)
			if path != tc.expectedPath {
				t.Fatalf("archive for %s did not match (%s != %s)", tc.docType.String(), path, tc.expectedPath)
			}
		})
	}

}
