package documents

import (
	"fmt"
	"os"

	"deepsolutionsvn.com/disc/providers/musicbrainz/archives"
)

const (
	MUSICBRAINZ_DROP_PATH    = "./MusicData/mb"
	MUSICBRAINZ_DROP_VERSION = "20220101"
	MUSICBRAINZ_STAGING_PATH = "./MusicData/mb/STAGING"
)

func CalcBuckets(mbid string) (string, string, string) {
	l1 := mbid[0:2]
	l2 := mbid[2:4]
	l3 := mbid[4:6]

	return l1, l2, l3
}

func MakeStagingPath(docType archives.DocumentType, stagingBasePath string, dropVersion string, did string) (stagingPath string, l1 string, l2 string, l3 string) {
	l1, l2, l3 = CalcBuckets(did)
	stagingPath = fmt.Sprintf("%s/%s/%s/%s/%s/%s", stagingBasePath, dropVersion, docType.ShortForm(), l1, l2, l3)

	return stagingPath, l1, l2, l3
}

func MakeContentFilename(stagingPath string, did string) (filename string) {
	filename = fmt.Sprintf("%s/%s.json", stagingPath, did)
	return filename
}

func WriteFileToStaging(docType archives.DocumentType, stagingBasePath string, dropVersion string, did string, content string) (filename string, l1 string, l2 string, l3 string, err error) {

	stagingPath, l1, l2, l3 := MakeStagingPath(docType, stagingBasePath, dropVersion, did)

	if _, err := os.Stat(stagingPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(stagingPath, 0755)
		} else {
			return "", "", "", "", err
		}
	}

	filename = MakeContentFilename(stagingPath, did)

	err = os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return "", "", "", "", err
	}

	return filename, l1, l2, l3, nil
}
