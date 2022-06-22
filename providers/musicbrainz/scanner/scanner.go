package scanner

import (
	"encoding/json"
	"fmt"
	"os"
)

type EntityInterface interface{}

type Entity struct {
	Value EntityInterface
	Error error
}

type JSONEntityStream struct {
	stream chan Entity
	count  int64
	closed bool
}

func NewJSONEntityStream() *JSONEntityStream {
	return &JSONEntityStream{
		stream: make(chan Entity, 6),
		count:  1,
		closed: false,
	}
}

func (es JSONEntityStream) Watch() <-chan Entity {
	return es.stream
}

func (es *JSONEntityStream) Start(entityFile string) {
	f, err := os.Open(entityFile)
	if err != nil {
		es.stream <- Entity{Error: fmt.Errorf("open file: %w", err)}
		return
	}
	defer f.Close()

	decoder := json.NewDecoder(f)

	es.count = int64(1)
	for decoder.More() {

		var e EntityInterface
		if err := decoder.Decode(&e); err != nil {
			es.stream <- Entity{Error: fmt.Errorf("decode line %d: %w", es.count, err)}
			return
		}

		es.stream <- Entity{
			Value: e,
		}

		es.count++
		// if es.count >= 10000 {
		// 	break
		// }
	}

	es.Close()
}

func (es *JSONEntityStream) Close() {
	if !es.closed {
		close(es.stream)
		es.closed = true
	}
}
