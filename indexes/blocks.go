package indexes

var BlockSize = int64(10000)

func GetBlockRange(count int64) (int64, int64) {
	start := (((count - 1) / BlockSize) * BlockSize) + 1
	end := start + BlockSize - 1
	return start, end
}
