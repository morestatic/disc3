package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"deepsolutionsvn.com/disc/ipfsutils"
)

func (opts *PublishMfsLibraryPathOptions) Execute(args []string) error {
	fmt.Printf("%#v\n", *opts)

	ipfsPath := opts.IpfsPath
	if ipfsPath == "" {
		ipfsPath = os.Getenv("IPFS_PATH")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("Spawning node on default repo")
	node, ipfs, err := ipfsutils.SpawnDefault(ctx, ipfsPath)
	if err != nil {
		return fmt.Errorf("failed to spawnDefault node: %w", err)
	}

	cid, err := ipfsutils.PublishMfsLibraryPath(ctx, node, ipfs, opts.LibraryIpnsRoot, opts.MutableFileStorePath)
	if err != nil {
		return fmt.Errorf("unable to push staging to ipfs: %w", err)
	}

	time.Sleep(1 * time.Second)

	// select {
	// case <-evtChan:
	// case <-ctx.Done():
	// }

	fmt.Printf("added library staging directory to IPFS with cid %s\n", cid)

	return err
}
