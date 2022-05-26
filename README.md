# DISC3

## Introduction

Disc3 takes publically available (CC0) music metadata and loads it into IPFS as a set of individual (approx 30M+) json files. Initially supporting metadata from discogs but more to come including musicbrainz and ISNI.

Note: This documentation is incomplete and very rough. Please see the code for more info. Docs will be fully updated in a future update.

## Setup Additional IPFS Node

- Check the Peer ID and top level CID at [www.disc3.xyz](https://www.disc3.xyz)
- Connect to the PeerID from your local IPFS node via `ipfs swarm connect <PeerID>`
- To keep a local copy of the data, add the CID via `ipfs add <cid>`, which could take a long time as a full copy of the metadata repo (approx. 150GB) will be made.

For a more stable replica node, sent a twitter DM to @morestatic with your PeerID and your peer will be added as a direct peer of the disc3 IPFS node.

## Setup Staging & Data Load

- Close the repo from here.
- Copy `env.base` to `.env` and update the environment variables to match your environment.

## Download Discogs XML Data

```bash
mkdir -p ~/discogs/data
wget http://discogs-data.s3-us-west-2.amazonaws.com/data/2022/discogs_20220401_releases.xml.gz
wget http://discogs-data.s3-us-west-2.amazonaws.com/data/2022/discogs_20220401_labels.xml.gz
wget http://discogs-data.s3-us-west-2.amazonaws.com/data/2022/discogs_20220401_masters.xml.gz
wget http://discogs-data.s3-us-west-2.amazonaws.com/data/2022/discogs_20220401_artists.xml.gz
```

## Staging Discogs Data

```bash
mkdir -p ~/db/sqlite
cat indexes/sql/sqlite_init_db.sql | sqlite3 ~/db/sqlite/indexes-20220401.db
mkdir -p ~/discogs/staging
go run disc3.go index discogs labels-xml
go run disc3.go index discogs masters-xml
go run disc3.go index discogs artists-xml
go run disc3.go index discogs releases-xml
go run disc3.go index discogs artists-releases
go run disc3.go stage discogs labels
go run disc3.go stage discogs masters
go run disc3.go stage discogs releases
go run disc3.go stage discogs artists
go run disc3.go stage discogs artists-releases
```

## Loading Data Into IPFS

```bash
ulimit -n 32768
export IPFS_PATH=~/ipfs
cd ~/staging
time find L/* -maxdepth 1 -regex '^L/[0-9][0-9]' -type d  -print0 | xargs -0 -n1 ~/dev/go/disc/cli/scripts/ipfs_add.sh
time find R/* -maxdepth 1 -regex '^R/[0-9][0-9]' -type d  -print0 | xargs -0 -n1 ~/dev/go/disc/cli/scripts/ipfs_add.sh
time find A/* -maxdepth 1 -regex '^A/[0-9][0-9]' -type d  -print0 | xargs -0 -n1 ~/dev/go/disc/cli/scripts/ipfs_add.sh
time find M/* -maxdepth 1 -regex '^M/[0-9][0-9]' -type d  -print0 | xargs -0 -n1 ~/dev/go/disc/cli/scripts/ipfs_add.sh

ipfs files mkdir -p /disc3/D/20220401
./cli/scripts/dump_artists_search_info.sh /home/ipfs/db/sqlite/indexes-20220401.db .
ipfs add artists_search_info.csv
ipfs files cp /ipfs/<cid-output-from-above> /disc3/D/20220401/artists_search_info.csv
ipfs add artists_search_info.csv.gz
ipfs files cp /ipfs/<cid-output-from-above> /disc3/D/20220401/artists_search_info.csv.gz

```

## TODO

## High-Level

- [ ] Improve recovery on IPFS node failure
- [ ] Automate the steps to load a new discogs drop to IPFS
- [ ] Add additional dedicated servers for IPFS nodes hosting the data
- [ ] Add Filecoin archiving of all data
- [ ] Improve documentation for data use + data loading
- [ ] Add support for musicbrainz.org metadata
- [ ] Implement a better website design (for disc3.xyz)
- [ ] Add a simple client side SDK to make the data even easily to load for developers
- [ ] Add additional indexes for metadata discovery
- [ ] Improve test code coverage
- [ ] General code improvements & refactoring

## Detailed

- [ ] Improve use of environment variables for scripts/data load
- [ ] Add tests for ipfs info service call
- [ ] Remove images element from staging json files
- [ ] Remove the underscore prefixes for converted xml attributes
- [ ] Add monitoring of the ipfs daemon
- [ ] Kill and restart ipfs daemon when using high memory
- [ ] Add option to run ipfs as docker container
- [ ] Integrate sharding options (code currently commented out)
- [ ] Better way to handle default cli options? no defaults? always require .env or options on the command line?
- [ ] Could goroutines speed up json staging?
- [ ] Use go code to upload docs from staging to ipfs (currently uses the IPFS cli)
- [ ] Either use or remove local ipfs code for uploads (see item above)
- [ ] Move db creation to golang
- [ ] Complete test code coverage
- [ ] Decide what to do about progress when adding docs to ipfs
- [ ] Think about whether a separate artist id is required? does this mean separate from the discogs id?
- [ ] Consider port of cli / command line args to cobra
