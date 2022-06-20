package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"deepsolutionsvn.com/disc/ipfsutils"
	archives "deepsolutionsvn.com/disc/providers/discogs/archives"
	documents "deepsolutionsvn.com/disc/providers/discogs/documents"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var ipfsReader ipfsutils.Reader = &ipfsutils.DefaultReader{}

type infoHandler func(c *gin.Context, dropId string, dt archives.DocumentType, useInfoSpec bool) (string, string, error)

type dispatchInfo struct {
	infoHandler infoHandler
	useInfoSpec bool
}

var dispatchTable = map[archives.DocumentType]dispatchInfo{
	archives.Artists:  {infoHandler: handleGetInfoJson, useInfoSpec: true},
	archives.Releases: {infoHandler: handleGetInfoJson, useInfoSpec: false},
	archives.Masters:  {infoHandler: handleGetInfoJson, useInfoSpec: false},
	archives.Labels:   {infoHandler: handleGetInfoJson, useInfoSpec: true},
}

func RepoServer(ipfsAddress string, port int64, dropId string) error {
	ipfsutils.DefaultIpfsAddress = ipfsAddress

	r := SetupRouter(port, nil, dropId)

	bindAddress := fmt.Sprintf("0.0.0.0:%s", strconv.FormatInt(port, 10))
	r.Run(bindAddress)

	return nil
}

func SetupRouter(port int64, reader ipfsutils.Reader, dropId string) *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Use(cors.Default())

	if reader != nil {
		ipfsReader = reader
	}

	r.GET("/ipfs/info", func(c *gin.Context) {
		mfsBasePath := fmt.Sprintf("%s/%s", ipfsutils.DefaultMfsPath, dropId)
		getIpfsInfo(c, mfsBasePath)
	})

	r.GET("/artist/discogs/:did", func(c *gin.Context) {
		dispatch(c, dropId, archives.Artists)
	})

	r.GET("/release/discogs/:did", func(c *gin.Context) {
		dispatch(c, dropId, archives.Releases)
	})

	r.GET("/master/discogs/:did", func(c *gin.Context) {
		dispatch(c, dropId, archives.Masters)
	})

	r.GET("/label/discogs/:did", func(c *gin.Context) {
		dispatch(c, dropId, archives.Labels)
	})

	return r
}

func getIpfsInfo(c *gin.Context, mfsBasePath string) {
	repoInfo, err := ipfsutils.GetRepoInfo(c, mfsBasePath)
	if err != nil {
		handleResponse(c, "", "", err)
		return
	}

	repoInfoJson, err := json.Marshal(repoInfo)
	if err != nil {
		handleResponse(c, "", "", err)
		return
	}

	handleResponse(c, string(repoInfoJson), "", err)
}

func dispatch(c *gin.Context, dropId string, dt archives.DocumentType) {
	documentJson, metaJson, err := getHandler(dt)(c, dropId, dt, getUseInfoSpec(dt))
	handleResponse(c, documentJson, metaJson, err)
}

func handleResponse(c *gin.Context, j string, m string, err error) {
	if err != nil {
		errString, _ := json.Marshal(err.Error())
		c.String(http.StatusBadGateway, fmt.Sprintf("{ \"message\": %s }", errString))
	} else {
		var response string
		if m == "" {
			response = fmt.Sprintf("{\n\"info\":%s\n}\n", j)
		} else {
			response = fmt.Sprintf("{\n\"meta\": %s,\n\"document\":%s\n}\n", m, j)
		}
		fmt.Println(response)
		c.String(http.StatusOK, response)
	}
}

func getHandler(dt archives.DocumentType) infoHandler {
	return dispatchTable[dt].infoHandler
}

func getUseInfoSpec(dt archives.DocumentType) bool {
	return dispatchTable[dt].useInfoSpec
}

func handleGetInfoJson(c *gin.Context, dropId string, dt archives.DocumentType, useInfoSpec bool) (string, string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	did := c.Param("did")
	log.Println(did)

	infoJson, metaJson, err := getInfoJson(ctx, dropId, did, useInfoSpec, dt)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(infoJson)
	}

	return infoJson, metaJson, err
}

func getInfoJson(ctx context.Context, dropId string, did string, useInfoSpec bool, dt archives.DocumentType) (string, string, error) {
	d, err := strconv.ParseInt(did, 10, 64)
	if err != nil {
		return "{}", "{}", err
	}

	basePath := fmt.Sprintf("%s/%s", ipfsutils.DefaultMfsPath, dropId)
	l1, l2, l3 := documents.CalcBuckets(d)

	bucketPath, err := documents.MakeBucketPathname(basePath, l1, l2, l3, dt.String())
	if err != nil {
		return "{}", "{}", err
	}

	var docSpec string
	if useInfoSpec {
		docSpec = fmt.Sprintf("%08d/%s_info.json", d, dt.Singular())
	} else {
		docSpec = fmt.Sprintf("%08d.json", d)
	}

	fileSpec := fmt.Sprintf("%s/%s", bucketPath, docSpec)

	log.Println(fileSpec)

	infoContent, contentInfo, err := ipfsReader.ReadAll(ctx, did, fileSpec)
	if err != nil {
		return "{}", "{}", err
	}

	contentInfo.Path = ipfsutils.PathInfo{
		MfsPath:       fileSpec,
		DiscogsPath:   ipfsutils.DefaultMfsPath,
		DocumentShard: dt.ShortForm(),
		LocationShard: fmt.Sprintf("%02d/%02d/%02d", l1, l2, l3),
		DocSpec:       docSpec,
	}
	metaJson, err := json.Marshal(contentInfo)
	if err != nil {
		return "{}", "{}", err
	}
	return infoContent, string(metaJson), nil
}
