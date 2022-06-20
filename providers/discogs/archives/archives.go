package archives

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strings"
)

type DiscogsScanProgress struct {
	Count    int64
	StartPos int64
	EndPos   int64
}

const (
	ArchiveDefaultPrefix    = "discogs"
	ArchiveDefaultSeparator = "_"
)

func MakeDefaultDiscogsArchiveName(archiveType DocumentType, dropVersion string, dropPath string) string {
	if dropPath == "" {
		dropPath = os.Getenv("DISCOGS_DROP_PATH")
	}
	if dropVersion == "" {
		dropVersion = os.Getenv("DISCOGS_DROP_VERSION")
	}
	name := fmt.Sprintf("%s/%s%s%s%s%s.xml",
		dropPath,
		ArchiveDefaultPrefix,
		ArchiveDefaultSeparator,
		dropVersion,
		ArchiveDefaultSeparator,
		archiveType)
	return name
}

func NewDiscogsFileReader(archiveType DocumentType, dropVersion string, dropPath string) (*os.File, error) {
	filename := MakeDefaultDiscogsArchiveName(archiveType, dropVersion, dropPath)

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func GetContent(f *os.File, startPos int64, endPos int64) (string, int64, error) {
	expectedCount := endPos - startPos
	bytes := make([]byte, expectedCount)

	n, err := f.ReadAt(bytes, startPos)
	readCount := int64(n)
	if err != nil {
		return "", readCount, err
	}

	if int64(n) != expectedCount {
		return "", readCount, errors.New("unexpected read length")
	}

	return strings.TrimSpace(string(bytes)), readCount, nil
}

func Parse(xmlContent string, o interface{}) error {
	err := xml.Unmarshal([]byte(xmlContent), o)
	fmt.Println(o)
	return err
}
