package server_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"deepsolutionsvn.com/disc/ipfsutils"
	"deepsolutionsvn.com/disc/server"
	"golang.org/x/net/context"
)

type TestIpfsReader struct{}

func (r *TestIpfsReader) ReadAll(ctx context.Context, did string, contentPath string) (string, *ipfsutils.ContentInfo, error) {
	fmt.Println(contentPath)
	contentInfo := &ipfsutils.ContentInfo{
		Cid: "Qm...",
	}
	if contentPath == "/disc3/D/20220401/A/00/00/00/00000020/artist_info.json" {
		content, err := ioutil.ReadFile("../testdata/staging/A/00/00/00/00000020/artist_info.json")
		return string(content), contentInfo, err
	}
	if contentPath == "/disc3/D/20220401/R/00/00/00/00000003.json" {
		content, err := ioutil.ReadFile("../testdata/staging/R/00/00/00/00000003.json")
		return string(content), contentInfo, err
	}
	return "{}", nil, nil
}

func TestGetArtistInfoJsonApi(t *testing.T) {
	router := server.SetupRouter(4000, &TestIpfsReader{}, "20220401")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/artist/discogs/20", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("unexpected response code (%d)", w.Code)
	}

	json := w.Body.String()
	if len(json) != 1875 {
		t.Fatalf("unexpected response length (%d)", len(json))
	}
}

func TestGetReleaseInfoJsonApi(t *testing.T) {
	router := server.SetupRouter(4000, &TestIpfsReader{}, "20220401")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/release/discogs/3", nil)
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("unexpected response code (%d)", w.Code)
	}

	json := w.Body.String()
	if len(json) != 15277 {
		t.Fatalf("unexpected response length (%d)", len(json))
	}
}
