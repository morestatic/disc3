package documents

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	archives "deepsolutionsvn.com/disc/providers/discogs/archives"
	xj "github.com/basgys/goxml2json"
)

type ReleaseDID int64
type ReleasesByArtist []ReleaseDID

func GetJson(xmlContent string) (string, error) {
	json, err := xj.Convert(strings.NewReader(string(xmlContent)), xj.WithAttrPrefix("_"), xj.WithContentPrefix(""))
	if err != nil {
		return "", err
	}
	return json.String(), nil
}

func AsJson(document interface{}) (string, error) {
	json, err := json.MarshalIndent(document, "", "    ")
	if err != nil {
		return "", err
	}

	return string(json), nil
}

func CalcBuckets(did int64) (int32, int32, int32) {
	l1 := did / 1000000
	l2 := (did - (l1 * 1000000)) / 10000
	l3 := (did - (l1 * 1000000) - (l2 * 10000)) / 100
	return int32(l1), int32(l2), int32(l3)
}

func GetStagingPath(stagingPath string) string {
	if stagingPath == "" {
		stagingPath = os.Getenv("STAGING_PATH")
	}
	return stagingPath
}

func GetShortNameForEntityType(entityType string) (string, error) {
	switch entityType {
	case archives.Artists:
		return "A", nil
	case archives.Releases:
		return "R", nil
	case archives.Labels:
		return "L", nil
	case archives.Masters:
		return "M", nil
	default:
		return "", errors.New("unknown archive type")
	}
}

func MakeBucketPathname(basePath string, l1 int32, l2 int32, l3 int32, entityType string) (string, error) {
	shortName, err := GetShortNameForEntityType(entityType)
	if err != nil {
		return "", err
	}
	path := fmt.Sprintf("%s/%s/%02d/%02d/%02d", basePath, shortName, l1, l2, l3)
	return path, err
}

func MakeContentFilename(bucketPath string, did int64, contentType string) string {
	filename := fmt.Sprintf("%s/%08d%s.json", bucketPath, did, contentType)
	return filename
}

func MakeArtistFolderName(stagingPath string, artistDid int64) (string, error) {
	l1, l2, l3 := CalcBuckets(artistDid)
	bucketPath, err := MakeBucketPathname(GetStagingPath(stagingPath), l1, l2, l3, archives.Artists)
	if err != nil {
		return "", err
	}
	artistFolderName := fmt.Sprintf("%s/%08d", bucketPath, artistDid)
	return artistFolderName, nil
}

func MakeLabelFolderName(stagingPath string, labelDid int64) (string, error) {
	l1, l2, l3 := CalcBuckets(labelDid)
	bucketPath, err := MakeBucketPathname(GetStagingPath(stagingPath), l1, l2, l3, archives.Labels)
	if err != nil {
		return "", err
	}
	labelFolderName := fmt.Sprintf("%s/%08d", bucketPath, labelDid)
	return labelFolderName, nil
}

func mkdir(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0755)
		} else {
			return err
		}
	}
	return err
}

func WriteArtistInfoToStaging(stagingPath string, did int64, content string) (string, error) {
	artistFolderName, err := MakeArtistFolderName(GetStagingPath(stagingPath), did)
	if err != nil {
		return "", err
	}

	err = mkdir(artistFolderName)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s/artist_info.json", artistFolderName)
	err = os.WriteFile(filename, []byte(content), 0644)

	return filename, err
}

func WriteLabelInfoToStaging(stagingPath string, did int64, content string) (string, error) {
	labelFolderName, err := MakeLabelFolderName(GetStagingPath(stagingPath), did)
	if err != nil {
		return "", err
	}

	err = mkdir(labelFolderName)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s/label_info.json", labelFolderName)
	err = os.WriteFile(filename, []byte(content), 0644)

	return filename, err
}

func MakeArtistReleasesFilename(stagingPath string, did int64) (string, error) {
	artistFolderName, err := MakeArtistFolderName(GetStagingPath(stagingPath), did)
	if err != nil {
		return "", err
	}

	err = mkdir(artistFolderName)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s/artist_releases.json", artistFolderName)
	return filename, nil
}

func WriteArtistReleasesToStaging(stagingPath string, did int64, content string) (string, error) {
	filename, err := MakeArtistReleasesFilename(GetStagingPath(stagingPath), did)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(filename, []byte(content), 0644)

	return filename, err
}

func ReadArtistReleasesJsonFromStaging(stagingPath string, did int64) (map[string]ReleasesByArtist, error) {
	filename, err := MakeArtistReleasesFilename(GetStagingPath(stagingPath), did)
	if err != nil {
		return nil, err
	}

	contentBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	releasesByArtist := make(map[string]ReleasesByArtist)
	err = json.Unmarshal(contentBytes, &releasesByArtist)
	if err != nil {
		return nil, err
	}

	return releasesByArtist, nil
}

func WriteFileToStaging(stagingPath string, did int64, content string, entityType string, contentType string) (string, int32, int32, int32, error) {
	l1, l2, l3 := CalcBuckets(did)
	bucketPath, err := MakeBucketPathname(GetStagingPath(stagingPath), l1, l2, l3, entityType)
	if err != nil {
		return "", 0, 0, 0, err
	}
	if _, err := os.Stat(bucketPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(bucketPath, 0755)
		} else {
			return "", 0, 0, 0, err
		}
	}
	filename := MakeContentFilename(bucketPath, did, contentType)
	err = os.WriteFile(filename, []byte(content), 0644)

	return filename, l1, l2, l3, err
}

func ArtistDocumentJsonExists(stagingPath string, did int64, docType string) (bool, error) {
	artistFolderName, err := MakeArtistFolderName(GetStagingPath(stagingPath), did)
	if err != nil {
		return false, err
	}

	filename := fmt.Sprintf("%s/artist_%s.json", artistFolderName, docType)
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func LabelDocumentJsonExists(stagingPath string, did int64, docType string) (bool, error) {
	labelFolderName, err := MakeLabelFolderName(GetStagingPath(stagingPath), did)
	if err != nil {
		return false, err
	}

	filename := fmt.Sprintf("%s/label_%s.json", labelFolderName, docType)
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func PrimaryDocumentExists(stagingPath string, entityType string, did int64, contentType string) (bool, error) {
	l1, l2, l3 := CalcBuckets(did)
	bucketPath, err := MakeBucketPathname(GetStagingPath(stagingPath), l1, l2, l3, entityType)
	if err != nil {
		return false, err
	}
	filename := MakeContentFilename(bucketPath, did, contentType)
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

// func PrimaryS2DocumentExists(stagingPath string, entityType string, did int64, contentType string) (bool, error) {
// 	l1, l2 := CalcS2Buckets(did)
// 	bucketPath, err := MakeS2BucketPathname(GetStagingPath(stagingPath), l1, l2, entityType)
// 	if err != nil {
// 		return false, err
// 	}
// 	filename := MakeContentFilename(bucketPath, did, contentType)
// 	if _, err := os.Stat(filename); err != nil {
// 		if os.IsNotExist(err) {
// 			return false, nil
// 		} else {
// 			return false, err
// 		}
// 	}
// 	return true, nil
// }

// func CalcS2Buckets(did int64) (int32, int32) {
// 	l1 := did / 1000000
// 	l2 := (did - (l1 * 1000000)) / 1000
// 	return int32(l1), int32(l2)
// }

// func MakeS2BucketPathname(basePath string, l1 int32, l2 int32, entityType string) (string, error) {
// 	shortName, err := GetShortNameForEntityType(entityType)
// 	if err != nil {
// 		return "", err
// 	}
// 	path := fmt.Sprintf("%s/%s/%03d/%03d", basePath, shortName, l1, l2)
// 	return path, err
// }

// func MakeLabelS2FolderName(stagingPath string, labelDid int64) (string, error) {
// 	l1, l2 := CalcS2Buckets(labelDid)
// 	bucketPath, err := MakeS2BucketPathname(GetStagingPath(stagingPath), l1, l2, archives.Labels)
// 	if err != nil {
// 		return "", err
// 	}
// 	labelFolderName := fmt.Sprintf("%s/%08d", bucketPath, labelDid)
// 	return labelFolderName, nil
// }

// func WriteS2FileToStaging(stagingPath string, did int64, content string, entityType string, contentType string) (string, int32, int32, error) {
// 	l1, l2 := CalcS2Buckets(did)
// 	bucketPath, err := MakeS2BucketPathname(GetStagingPath(stagingPath), l1, l2, entityType)
// 	if err != nil {
// 		return "", 0, 0, err
// 	}
// 	if _, err := os.Stat(bucketPath); err != nil {
// 		if os.IsNotExist(err) {
// 			os.MkdirAll(bucketPath, 0755)
// 		} else {
// 			return "", 0, 0, err
// 		}
// 	}
// 	filename := MakeContentFilename(bucketPath, did, contentType)
// 	err = os.WriteFile(filename, []byte(content), 0644)
// 	return filename, l1, l2, err
// }

// func ArtistS2DocumentJsonExists(stagingPath string, did int64, docType string) (bool, error) {
// 	artistFolderName, err := MakeArtistS2FolderName(GetStagingPath(stagingPath), did)
// 	if err != nil {
// 		return false, err
// 	}
// 	filename := fmt.Sprintf("%s/artist_%s.json", artistFolderName, docType)
// 	if _, err := os.Stat(filename); err != nil {
// 		if os.IsNotExist(err) {
// 			return false, nil
// 		} else {
// 			return false, err
// 		}
// 	}
// 	return true, nil
// }

// func MakeArtistS2FolderName(stagingPath string, artistDid int64) (string, error) {
// 	l1, l2 := CalcS2Buckets(artistDid)
// 	bucketPath, err := MakeS2BucketPathname(GetStagingPath(stagingPath), l1, l2, archives.Artists)
// 	if err != nil {
// 		return "", err
// 	}
// 	artistFolderName := fmt.Sprintf("%s/%08d", bucketPath, artistDid)
// 	return artistFolderName, nil
// }

// func LabelS2DocumentJsonExists(stagingPath string, did int64, docType string) (bool, error) {
// 	labelFolderName, err := MakeLabelS2FolderName(GetStagingPath(stagingPath), did)
// 	if err != nil {
// 		return false, err
// 	}
// 	filename := fmt.Sprintf("%s/label_%s.json", labelFolderName, docType)
// 	if _, err := os.Stat(filename); err != nil {
// 		if os.IsNotExist(err) {
// 			return false, nil
// 		} else {
// 			return false, err
// 		}
// 	}
// 	return true, nil
// }
