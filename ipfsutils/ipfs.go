package ipfsutils

import (
	"context"
	"errors"
	"fmt"

	gopath "path"

	"deepsolutionsvn.com/disc/archives"
	"deepsolutionsvn.com/disc/documents"

	// "deepsolutionsvn.com/disc/ipfsutils"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-ipfs/core"
	mfs "github.com/ipfs/go-mfs"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/ipfs/interface-go-ipfs-core/path"
)

var DefaultIpnsPath = "/ipns/k51qzi5uqu5dgmon3ymmr29hi09r82jskuk8dec8fp8fs3yipr3201misx22tb"
var DefaultMfsPath = "/disc3/D"

func GetIPNSPath(ipnsPath string) string {
	if ipnsPath == "" {
		ipnsPath = DefaultIpnsPath
	}
	return ipnsPath
}

func GetMfsPath(mfsPath string) string {
	if mfsPath == "" {
		mfsPath = DefaultMfsPath
	}
	return mfsPath
}

func GetArtistFolderCid(ctx context.Context, ipfs icore.CoreAPI, ipnsPath string, artistDid int64) (string, cid.Cid, error) {
	l1, l2, l3 := documents.CalcBuckets(artistDid)
	bucketPath, err := documents.MakeBucketPathname(GetIPNSPath(ipnsPath), l1, l2, l3, archives.Artists)
	if err != nil {
		return "", cid.Cid{}, err
	}
	folderPath := fmt.Sprintf("%s/%08d", bucketPath, artistDid)

	path, err := ipfs.ResolvePath(ctx, path.New(folderPath))
	if err != nil {
		return "", cid.Cid{}, err
	}
	cid := path.Cid()

	return folderPath, cid, nil
}

func GetReleaseJsonCid(ctx context.Context, ipfs icore.CoreAPI, ipnsPath string, releaseDid int64) (string, cid.Cid, error) {
	l1, l2, l3 := documents.CalcBuckets(releaseDid)
	bucketPath, err := documents.MakeBucketPathname(GetIPNSPath(ipnsPath), l1, l2, l3, archives.Releases)
	if err != nil {
		return "", cid.Cid{}, err
	}
	folderPath := fmt.Sprintf("%s/%08d.json", bucketPath, releaseDid)

	path, err := ipfs.ResolvePath(ctx, path.New(folderPath))
	if err != nil {
		return "", cid.Cid{}, err
	}
	cid := path.Cid()

	return folderPath, cid, nil
}

func printAddEvent(evt interface{}) {
	info := evt.(*icore.AddEvent)
	fmt.Printf("%s -> %s\n", info.Name, info.Path.Cid())
}

func drainAddEvents(evtChan chan interface{}) {
	for len(evtChan) != 0 {
		evt := <-evtChan
		printAddEvent(evt)
	}
}

func setupProgressMonitor(ctx context.Context) (chan interface{}, error) {
	evtChan := make(chan interface{}, 8)
	done := false
	go func() {
		for !done {
			select {
			case evt := <-evtChan:
				printAddEvent(evt)
			case <-ctx.Done():
				drainAddEvents(evtChan)
				done = true
			}
		}
	}()
	return evtChan, nil
}

func PushStagingToIpfs(ctx context.Context, node *core.IpfsNode, ipfs icore.CoreAPI, stagingPath string, noMfs bool, mfsPath string, replaceMfs bool) (string, error) {
	stagingFolder, err := GetUnixfsNode(documents.GetStagingPath(stagingPath))
	if err != nil {
		return "", fmt.Errorf("could not get staging folder: %w", err)
	}

	evtChan, err := setupProgressMonitor(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to setup progress handler: %w", err)
	}

	opts := []options.UnixfsAddOption{}
	opts = append(opts, options.Unixfs.Events(evtChan))
	// opts = append(opts, options.Unixfs.Nocopy(true))

	stagedDirectoryCid, err := ipfs.Unixfs().Add(ctx, stagingFolder, opts...)
	if err != nil {
		return "", fmt.Errorf("could not add directory: %w", err)
	}

	if !noMfs {
		sourceCid := "/ipfs/" + stagedDirectoryCid.Cid().String()
		destPath := GetMfsPath(mfsPath)
		fmt.Printf("copying %s -> %s\n", sourceCid, destPath)

		parentPath, name := gopath.Split(destPath)

		parentNode, err := getDirNode(node.FilesRoot, parentPath)
		if err != nil {
			newDirOpts := mfs.MkdirOpts{
				Mkparents: true,
			}
			err = mfs.Mkdir(node.FilesRoot, parentPath, newDirOpts)
			if err != nil {
				return "", fmt.Errorf("unable to create parent node %s, err: %w", parentPath, err)
			}
		} else {
			if replaceMfs {
				err = parentNode.Unlink(name)
				if err != nil {
					return "", fmt.Errorf("unable to remove %s, err: %w", name, err)
				}
			}
		}

		sourceNode, err := ipfs.ResolveNode(ctx, path.New(sourceCid))
		if err != nil {
			return "", errors.New("staged cid not found")
		}

		err = mfs.PutNode(node.FilesRoot, destPath, sourceNode)
		if err != nil {
			return "", fmt.Errorf("unable to copy staged directory to mutable filesystem destination: %w", err)
		}

		_, err = mfs.FlushPath(ctx, node.FilesRoot, DefaultMfsPath)
		if err != nil {
			return "", fmt.Errorf("unable to flush destination: %w", err)
		}
	}

	return stagedDirectoryCid.String(), nil
}

func getDirNode(root *mfs.Root, path string) (*mfs.Directory, error) {
	dirNode, err := mfs.Lookup(root, path)
	if err != nil {
		return nil, err
	}

	directory, ok := dirNode.(*mfs.Directory)
	if !ok {
		return nil, errors.New("expected *mfs.Directory, didn't get it. This is likely a race condition")
	}
	return directory, nil
}

func PublishMfsLibraryPath(ctx context.Context, node *core.IpfsNode, ipfs icore.CoreAPI, ipnsPath string, mfsPath string) (string, error) {

	mfsNode, err := mfs.Lookup(node.FilesRoot, GetMfsPath(mfsPath))
	if err != nil {
		return "", fmt.Errorf("unable to lookup mfs path: %s, err: %w", mfsPath, err)
	}
	nodeDetails, err := mfsNode.GetNode()
	if err != nil {
		return "", fmt.Errorf("unable to lookup get node for mfs path: %s, err: %w", mfsPath, err)
	}

	cid := nodeDetails.Cid()
	sourceCidPath := path.New("/ipfs/" + cid.String())

	// targetIpnsPath := GetIPNSPath(ipnsPath)

	// if verifyExists, _ := req.Options[resolveOptionName].(bool); verifyExists {
	// 	_, err := ipfs.ResolveNode(req.Context, p)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	opts := []options.NamePublishOption{
		options.Name.AllowOffline(true),
		options.Name.Key("self"),
	}

	out, err := ipfs.Name().Publish(ctx, sourceCidPath, opts...)
	if err != nil {
		// if err == icore.ErrOffline {
		// 	return fmt.Errorf("")
		// 	err = errAllowOffline
		// }
		return "", fmt.Errorf("unable to publish %s, err: %w", mfsPath, err)
	}

	fmt.Println(out.Name())

	return sourceCidPath.String(), nil
}
