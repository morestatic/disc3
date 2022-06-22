package archives

import "fmt"

func MakeDefaultMusicBrainzArchiveName(dt DocumentType, dropBasePath string, dropVersion string) (path string) {

	path = fmt.Sprintf("%s/%s/%s/mbdump/%s.json", dropBasePath, dropVersion, dt.ShortForm(), dt.ArchiveBaseName())

	return path
}
