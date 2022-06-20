package archives

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"reflect"
)

type DiscogsFileDecoder struct {
	xd *xml.Decoder
	f  *os.File
}

type DiscogsElementProcessor func(*DiscogsFileDecoder, interface{}, interface{}, DiscogsScanProgress) (bool, error)

func NewDiscogsFileDecoder(archiveType DocumentType, dropVersion string, dropPath string) (*DiscogsFileDecoder, error) {
	filename := MakeDefaultDiscogsArchiveName(archiveType, dropVersion, dropPath)

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	decoder := xml.NewDecoder(f)

	d := &DiscogsFileDecoder{
		f:  f,
		xd: decoder,
	}

	return d, nil
}

func (d DiscogsFileDecoder) GetFile() *os.File {
	return d.f
}

func (d *DiscogsFileDecoder) Close() {
	d.f.Close()
}

func (d *DiscogsFileDecoder) Scan(dt DocumentType, elementName string, outputSample interface{}, processorFn DiscogsElementProcessor) (bool, int64, error) {
	var start int64 = d.xd.InputOffset()
	var end int64 = 0
	var quit = false
	var err error = nil
	var t xml.Token = nil
	var count int64 = 0

	for {
		// start := decoder.InputOffset()
		t, err = d.xd.Token()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return true, count, fmt.Errorf("unable to decode token: %w", err)
			}
		}

		if t == nil {
			break
		}

		switch e := t.(type) {
		case xml.StartElement:
			// fmt.Println(e)
			if e.Name.Local == string(dt) {
				start = d.xd.InputOffset()
			}
			if e.Name.Local == elementName {
				if count != 0 {
					start = end
				}

				count++

				outputType := reflect.TypeOf(outputSample)
				outputValue := reflect.New(outputType)
				o := outputValue.Interface()

				d.xd.DecodeElement(o, &e)
				end = d.xd.InputOffset()

				quit, err = processorFn(d, &e, o, DiscogsScanProgress{Count: count, StartPos: start, EndPos: end})
			}
		}

		if quit || err != nil {
			break
		}
	}
	return quit, count, err
}
