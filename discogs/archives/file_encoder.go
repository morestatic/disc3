package archives

import (
	"errors"
	"fmt"
	"os"
)

type DiscogsFileEncoder struct {
	f           *os.File
	archiveType DocumentType
}

func NewDiscogsFileEncoder(archiveType DocumentType, filename string) (*DiscogsFileEncoder, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	e := &DiscogsFileEncoder{
		f:           f,
		archiveType: archiveType,
	}

	startElement := fmt.Sprintf("<%s>", archiveType)
	f.WriteString(startElement)

	return e, nil
}

func (e *DiscogsFileEncoder) Write(contents string) error {
	bytesWritten, err := e.f.Write([]byte(contents))
	if err != nil {
		return err
	}
	if bytesWritten != len(contents) {
		return errors.New("unexpected number of bytes written")
	}
	return nil
}

func (e *DiscogsFileEncoder) Close() {
	endElement := fmt.Sprintf("</%s>\n", e.archiveType)
	e.f.WriteString(endElement)
	e.f.Close()
}
