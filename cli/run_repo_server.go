package cli

import (
	"fmt"
	"os"
	"strconv"

	"deepsolutionsvn.com/disc/server"
)

func (opts *RepoServerOptions) Execute(args []string) error {
	var err error

	fmt.Printf("%#v\n", *opts)

	port := opts.Port
	if port == 0 {
		portStr := os.Getenv("PORT")
		port, err = strconv.ParseInt(portStr, 10, 0)
		if err != nil {
			return fmt.Errorf("cannot parse port option (%s)", portStr)
		}
	}

	address := os.Getenv("IPFS_ADDRESS")
	if address == "" {
		address = opts.IpfsAddress
	}

	dropId := os.Getenv("DISCOGS_DROP_VERSION")
	if dropId == "" {
		dropId = opts.DropId
	}

	fmt.Println(address)
	server.RepoServer(address, port, dropId)

	return nil
}
