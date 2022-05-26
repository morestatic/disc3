package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"deepsolutionsvn.com/disc/ipfsutils"
)

func (opts *PushToIpfsOptions) Execute(args []string) error {
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

	cid, err := ipfsutils.PushStagingToIpfs(ctx, node, ipfs, opts.LibStagingPath, opts.NoMutableFileStore, opts.MutableFileStorePath, opts.Replace)
	if err != nil {
		return fmt.Errorf("unable to push staging to ipfs: %w", err)
	}

	fmt.Println("DING")

	time.Sleep(1 * time.Second)

	// select {
	// case <-evtChan:
	// case <-ctx.Done():
	// }

	fmt.Printf("added library staging directory to IPFS with cid %s\n", cid)

	return err
}
