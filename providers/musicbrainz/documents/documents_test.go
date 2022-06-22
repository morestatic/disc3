package documents_test

import (
	"testing"

	"deepsolutionsvn.com/disc/providers/musicbrainz/archives"
	"deepsolutionsvn.com/disc/providers/musicbrainz/documents"
)

func TestShouldCalculateStagingPath(t *testing.T) {
	cases := []struct {
		testName        string
		docType         archives.DocumentType
		stagingBasePath string
		dropVersion     string
		did             string
		expectedPath    string
	}{
		{
			testName:        "Artist",
			docType:         archives.Artists,
			stagingBasePath: ".",
			dropVersion:     "20220101",
			did:             "470a4ced-1323-4c91-8fd5-0bb3fb4c932a",
			expectedPath:    "./20220101/A/47/0a/4c",
		},
		{
			testName:        "Master",
			docType:         archives.Masters,
			stagingBasePath: ".",
			dropVersion:     "20220101",
			did:             "b23eb64a-51aa-3d95-8c44-bc78af7b0457",
			expectedPath:    "./20220101/M/b2/3e/b6",
		},
		{
			testName:        "Release",
			docType:         archives.Releases,
			stagingBasePath: ".",
			dropVersion:     "20220102",
			did:             "ca2aad6f-f0ac-417b-b485-9df940c8fa48",
			expectedPath:    "./20220102/R/ca/2a/ad",
		},
		{
			testName:        "Label",
			docType:         archives.Labels,
			stagingBasePath: ".",
			dropVersion:     "20220102",
			did:             "888875e6-a416-494c-91fe-0cf35295f34f",
			expectedPath:    "./20220102/L/88/88/75",
		},
	}

	for _, tc := range cases {
		t.Run(tc.testName, func(t *testing.T) {
			stagingPath, _, _, _ := documents.MakeStagingPath(tc.docType, tc.stagingBasePath, tc.dropVersion, tc.did)
			if stagingPath != tc.expectedPath {
				t.Fatalf("stagingPath for %s did not match (%s != %s)", tc.docType.String(), stagingPath, tc.expectedPath)
			}
		})
	}
}

func TestShouldMakeContentFilename(t *testing.T) {
	contentFilename := documents.MakeContentFilename("1", "2")
	expectedFilename := "1/2.json"

	if contentFilename != expectedFilename {
		t.Fatalf("filename did not match (%s != %s)", contentFilename, expectedFilename)
	}
}
