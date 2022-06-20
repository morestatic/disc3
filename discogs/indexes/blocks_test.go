package indexes_test

import (
	"testing"

	indexes "deepsolutionsvn.com/disc/discogs/indexes"
)

func TestShouldGetFirstBlockRange(t *testing.T) {
	count := int64(1)
	start, end := indexes.GetBlockRange(count)
	expectedStart := int64(1)
	expectedEnd := int64(10000)

	if start != expectedStart {
		t.Fatal("unexpected start offset for block", start, expectedStart)
	}
	if end != expectedEnd {
		t.Fatal("unexpected end offset for block", end, expectedEnd)
	}
}

func TestShouldGetSecondBlockRange(t *testing.T) {
	blockNumber := int64(2)
	count := indexes.BlockSize * blockNumber
	start, end := indexes.GetBlockRange(count)
	expectedStart := int64(10001)
	expectedEnd := int64(20000)
	if start != expectedStart {
		t.Fatal("unexpected start offset for block", start, expectedStart)
	}
	if end != expectedEnd {
		t.Fatal("unexpected end offset for block", end, expectedEnd)
	}
}
