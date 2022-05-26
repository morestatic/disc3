package ipfsutils

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	shell "github.com/ipfs/go-ipfs-api"

	files "github.com/ipfs/go-ipfs-files"
	icore "github.com/ipfs/interface-go-ipfs-core"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

type Reader interface {
	ReadAll(ctx context.Context, did string, contentPath string) (string, *ContentInfo, error)
}

type DefaultReader struct{}

type PathInfo struct {
	MfsPath       string `json:"mfsPath"`
	DiscogsPath   string `json:"discogsPath"`
	DocumentShard string `json:"documentShard"`
	LocationShard string `json:"locationShard"`
	DocSpec       string `json:"docSpec"`
}

type ContentInfo struct {
	Cid  string   `json:"cid"`
	Did  string   `json:"did"`
	Path PathInfo `json:"path"`
}

type RepoInfo struct {
	Cid    string `json:"cid"`
	PeerID string `json:"peerID"`
}

var DefaultIpfsAddress = "127.0.0.1:5001"

// var DefaultIpfsAddress = "host.docker.internal:5001"

func (r *DefaultReader) ReadAll(ctx context.Context, did string, contentPath string) (string, *ContentInfo, error) {
	fmt.Printf("contentPath = %s\n", contentPath)

	sh := shell.NewShell(DefaultIpfsAddress)
	reader, err := sh.FilesRead(ctx, contentPath)
	if err != nil {
		return "", nil, err
	}
	defer reader.Close()

	info, err := sh.FilesStat(ctx, contentPath)
	if err == nil {
		log.Printf("%+v", info)
	}

	contentInfo := &ContentInfo{
		Cid: info.Hash,
		Did: did,
	}

	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", nil, err
	}

	return string(buf), contentInfo, nil
}

func GetRepoInfo(ctx context.Context, mfsBase string) (repoInfo *RepoInfo, err error) {
	fmt.Printf("mfsBase = %s\n", mfsBase)

	sh := shell.NewShell(DefaultIpfsAddress)

	dirInfo, err := sh.FilesStat(ctx, mfsBase)
	if err != nil {
		return nil, err
	}

	nodeInfo, err := sh.ID()
	if err != nil {
		return nil, err
	}

	repoInfo = &RepoInfo{
		Cid:    dirInfo.Hash,
		PeerID: nodeInfo.ID,
	}

	return repoInfo, nil
}

func GetIpfsPath(ipfsPath string) string {
	if ipfsPath == "" {
		ipfsPath = os.Getenv("IPFS_PATH")
	}
	return ipfsPath
}

func CreateNode(ctx context.Context, repoPath string) (*core.IpfsNode, icore.CoreAPI, error) {
	// Open the repo
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, nil, err
	}

	// Construct the node
	nodeOptions := &core.BuildCfg{
		Online: true,
		// Routing: libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
		Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		Repo:    repo,
	}

	node, err := core.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, nil, err
	}

	// Attach the Core API to the constructed node
	api, err := coreapi.NewCoreAPI(node)
	return node, api, err
}

func SpawnDefault(ctx context.Context, ipfsPath string) (*core.IpfsNode, icore.CoreAPI, error) {
	defaultPath := ipfsPath

	if err := setupPlugins(defaultPath); err != nil {
		return nil, nil, err
	}

	return CreateNode(ctx, defaultPath)
}

func setupPlugins(externalPluginsPath string) error {
	// Load any external plugins if available on externalPluginsPath
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	// Load preloaded and external plugins
	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}

func GetUnixfsNode(path string) (files.Node, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	f, err := files.NewSerialFile(path, false, st)
	if err != nil {
		return nil, err
	}

	return f, nil
}
