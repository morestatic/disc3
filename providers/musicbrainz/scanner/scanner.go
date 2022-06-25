package scanner

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

const (
	EntityChannelBufferSize = 12
)

type EntityOpenFileErr struct {
	Err error
}

func (e EntityOpenFileErr) Error() string {
	return e.Err.Error()
}

type EntityStreamDecodeErr struct {
	Count int64
	Err   error
}

func (e EntityStreamDecodeErr) Error() string {
	return fmt.Sprintf("decode error at entity %d: %v", e.Count, e.Err.Error())
}

type EntityInterface interface{}

type Entity struct {
	Value EntityInterface
	Error error
}

type JSONEntityStream struct {
	stream chan Entity
	count  int64
	mut    sync.RWMutex
}

func NewJSONEntityStream() *JSONEntityStream {
	return &JSONEntityStream{
		stream: make(chan Entity, EntityChannelBufferSize),
		count:  1,
	}
}

func (es *JSONEntityStream) Watch() <-chan Entity {
	return es.stream
}

func (es *JSONEntityStream) GetCount() int64 {
	es.mut.RLock()
	defer es.mut.RUnlock()
	return es.count
}

func (es *JSONEntityStream) Start(entityFile string) {
	f, err := os.Open(entityFile)
	if err != nil {
		es.stream <- Entity{Error: EntityOpenFileErr{Err: err}}
		return
	}
	defer f.Close()

	decoder := json.NewDecoder(f)

	for decoder.More() {

		var e EntityInterface
		if err := decoder.Decode(&e); err != nil {
			es.stream <- Entity{Error: &EntityStreamDecodeErr{Count: es.count, Err: err}}
			return
		}

		es.stream <- Entity{
			Value: e,
		}

		es.mut.Lock()
		es.count++
		// if es.count >= 1000 {
		// 	es.mut.Unlock()
		// 	es.stream <- Entity{Error: fmt.Errorf("break %d", es.count)}
		// 	close(es.stream)
		// 	break
		// }
		es.mut.Unlock()
	}

	close(es.stream)
}
