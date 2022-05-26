package cli

import (
	"context"
	"fmt"
	"strconv"

	"deepsolutionsvn.com/disc/ipfsutils"
)

func (opts *GetArtistFolderCidOptions) Execute(args []string) error {

	did, err := strconv.ParseInt(opts.DiscId, 10, 64)
	if err != nil {
		return fmt.Errorf("unable to parse did: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("Spawning node on default repo")
	_, ipfs, err := ipfsutils.SpawnDefault(ctx, opts.IpfsPath)
	if err != nil {
		return fmt.Errorf("failed to spawnDefault node: %w", err)
	}

	fmt.Println("Fetching cid for artist folder")
	folderPath, cid, err := ipfsutils.GetArtistFolderCid(ctx, ipfs, opts.LibraryIpnsRoot, did)
	if err != nil {
		return fmt.Errorf("failed to get artist folder cid: %w", err)
	}

	fmt.Println(folderPath)
	fmt.Println(cid)

	return err
}
