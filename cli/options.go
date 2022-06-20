package cli

import (
	"os"
	"os/signal"
	"syscall"
)

type ArchiveOptions struct {
	DropId   string `short:"d" long:"drop" description:"music metadata drop id (usually date)"`
	DropPath string `short:"p" long:"path" description:"music metadata drop path"`
}

type IndexOptions struct {
	ConnString string `short:"c" long:"conn" description:"database connection string"`
	StoreType  string `short:"i" long:"index-type" description:"store type (default: sqlite)" default:"sqlite"`
}

type StagingOptions struct {
	LibStagingPath string `short:"l" long:"lib" description:"music metadata json document staging path"`
}

type GeneralOptions struct {
	Quiet bool `bool:"q" long:"quiet" description:"quiet mode, do not show progress meter"`
}

type CreateArtistsXmlIndexOptions struct {
	ArchiveOptions
	IndexOptions
	GeneralOptions
}

type CreateLabelsXmlIndexOptions struct {
	ArchiveOptions
	IndexOptions
	GeneralOptions
}

type CreateReleasesXmlIndexOptions struct {
	ArchiveOptions
	IndexOptions
	GeneralOptions
}

type CreateMastersXmlIndexOptions struct {
	ArchiveOptions
	IndexOptions
	GeneralOptions
}

type CreateArtistsReleasesIndexOptions struct {
	ArchiveOptions
	IndexOptions
	GeneralOptions
}

type CreateArtistsSearchInfoOptions struct {
	ArchiveOptions
	IndexOptions
	GeneralOptions
}

type CreateReleaseGenresListOptions struct {
	ArchiveOptions
	IndexOptions
	GeneralOptions
}

type CreateReleaseStylesListOptions struct {
	ArchiveOptions
	IndexOptions
	GeneralOptions
}

type MakeArtistsXmlArchiveOptions struct {
	OutputName string `short:"o" long:"output" description:"output filename (required)" required:"true"`
	ArchiveOptions
	IndexOptions
	StagingOptions
}

type MakeArtistsReleasesXmlArchiveOptions struct {
	OutputName string `short:"o" long:"output" description:"output filename (required)" required:"true"`
	ArchiveOptions
	IndexOptions
	StagingOptions
}

type ExtractArtistJsonOptions struct {
	DiscId string `short:"a" long:"artist" description:"artist id (required)" required:"true"`
	ArchiveOptions
	IndexOptions
}

type ExtractReleaseJsonOptions struct {
	DiscId string `short:"r" long:"release" description:"release id (required)" required:"true"`
	ArchiveOptions
	IndexOptions
}

type DiscogsPutArtistJsonOptions struct {
	DiscId string `short:"a" long:"did" description:"artist id (required)" required:"true"`
	ArchiveOptions
	IndexOptions
	StagingOptions
}

type DiscogsPutArtistsJsonOptions struct {
	ArchiveOptions
	IndexOptions
	GeneralOptions
	StagingOptions
}

type MusicBrainzPutArtistsJsonOptions struct {
	ArchiveOptions
	IndexOptions
	GeneralOptions
	StagingOptions
}

type DiscogsPutArtistsReleasesJsonOptions struct {
	ArchiveOptions
	IndexOptions
	GeneralOptions
	StagingOptions
}

type MusicBrainzPutArtistsReleasesJsonOptions struct {
	ArchiveOptions
	IndexOptions
	GeneralOptions
	StagingOptions
}

type DiscogsPutLabelsJsonOptions struct {
	ArchiveOptions
	IndexOptions
	GeneralOptions
	StagingOptions
}

type MusicBrainzPutLabelsJsonOptions struct {
	ArchiveOptions
	IndexOptions
	GeneralOptions
	StagingOptions
}

type DiscogsPutReleaseJsonOptions struct {
	DiscId string `short:"r" long:"did" description:"release id (required)" required:"true"`
	ArchiveOptions
	IndexOptions
	GeneralOptions
	StagingOptions
}

type DiscogsPutReleasesJsonOptions struct {
	ArchiveOptions
	IndexOptions
	GeneralOptions
	StagingOptions
}

type MusicBrainzPutReleasesJsonOptions struct {
	ArchiveOptions
	IndexOptions
	GeneralOptions
	StagingOptions
}

type DiscogsPutMastersJsonOptions struct {
	ArchiveOptions
	IndexOptions
	GeneralOptions
	StagingOptions
}

type MusicBrainzPutMastersJsonOptions struct {
	ArchiveOptions
	IndexOptions
	GeneralOptions
	StagingOptions
}

type DiscogsPutArtistReleasesJsonOptions struct {
	DiscId string `short:"a" long:"did" description:"artist id (required)" required:"true"`
	ArchiveOptions
	IndexOptions
	GeneralOptions
	StagingOptions
}

type PushToIpfsOptions struct {
	IpfsPath             string `short:"i" long:"ipfs" description:"ipfs path"`
	LibraryIpnsRoot      string `short:"n" long:"ipns" description:"ipns base path"`
	LibStagingPath       string `short:"l" long:"lib" description:"discogs json document staging path"`
	MutableFileStorePath string `short:"o" long:"output" description:"output path on the ipfs mutable filesystem"`
	Replace              bool   `short:"r" long:"replace" decription:"replace existing mutable filesystem location"`
	NoMutableFileStore   bool   `long:"no-mfs" decription:"create the cid only, no use of the mutable filesystem"`
}

type PublishMfsLibraryPathOptions struct {
	IpfsPath             string `short:"i" long:"ipfs" description:"ipfs path"`
	LibraryIpnsRoot      string `short:"n" long:"ipns" description:"ipns base path"`
	MutableFileStorePath string `short:"o" long:"output" description:"output path on the ipfs mutable filesystem"`
}

type GetArtistFolderCidOptions struct {
	DiscId          string `short:"d" long:"did" description:"artist id (required)" required:"true"`
	IpfsPath        string `short:"i" long:"ipfs" description:"ipfs path"`
	LibraryIpnsRoot string `short:"n" long:"ipns" description:"ipns base path"`
}

type GetReleaseJsonCidOptions struct {
	DiscId          string `short:"d" long:"did" description:"release id (required)" required:"true"`
	IpfsPath        string `short:"i" long:"ipfs" description:"ipfs path"`
	LibraryIpnsRoot string `short:"n" long:"ipns" description:"ipns base path"`
}

type RepoServerOptions struct {
	IpfsAddress          string `short:"I" long:"ipfs" description:"ipfs address" default:"127.0.0.1:5001"`
	Port                 int64  `short:"P" long:"port" description:"server port" default:"4000"`
	DropId               string `short:"d" long:"drop" description:"discogs drop (usually date)"`
	MutableFileStorePath string `short:"o" long:"output" description:"output path on the ipfs mutable filesystem"`
}

type DiscogsIndexCommands struct {
	CreateArtistsXmlIndex        CreateArtistsXmlIndexOptions      `command:"artists-xml" description:"Create xml index of the artists archive" required:"true"`
	CreateLabelsXmlIndex         CreateLabelsXmlIndexOptions       `command:"labels-xml" description:"Create xml index of the labels archive" required:"true"`
	CreateReleasesXmlIndex       CreateReleasesXmlIndexOptions     `command:"releases-xml" description:"Create xml index of the releases archive" required:"true"`
	CreateMastersXmlIndex        CreateMastersXmlIndexOptions      `command:"masters-xml" description:"Create xml index of the masters archive" required:"true"`
	CreateArtistsReleasesIndexes CreateArtistsReleasesIndexOptions `command:"artists-releases" description:"Create artist and label indexes of releases" required:"true"`
	CreateArtistsSearchInfo      CreateArtistsSearchInfoOptions    `command:"artists-search-info" description:"Create artist search info" required:"true"`
	CreateReleaseGenresList      CreateReleaseGenresListOptions    `command:"release-genres-list" description:"Create list of genres used by releases" required:"true"`
	CreateReleaseStylesList      CreateReleaseStylesListOptions    `command:"release-styles-list" description:"Create list of styles used by releases" required:"true"`
}

type DiscogsStagingCommands struct {
	DiscogsPutArtistJson          DiscogsPutArtistJsonOptions          `command:"artist" description:"Put json from artist archive in json staging" required:"true"`
	DiscogsPutArtistReleasesJson  DiscogsPutArtistReleasesJsonOptions  `command:"artist-releases" description:"Put json from releases archive in json staging for the artist" required:"true"`
	DiscogsPutArtistsJson         DiscogsPutArtistsJsonOptions         `command:"artists" description:"Put json for all artists from the artists archive in json staging" required:"true"`
	DiscogsPutArtistsReleasesJson DiscogsPutArtistsReleasesJsonOptions `command:"artists-releases" description:"Put artist releases json in json staging for all artists" required:"true"`
	DiscogsPutLabelsJson          DiscogsPutLabelsJsonOptions          `command:"labels" description:"Put json for all labels from the labels archive in json staging" required:"true"`
	DiscogsPutReleaseJson         DiscogsPutReleaseJsonOptions         `command:"release" description:"Put json from release archive in json staging" required:"true"`
	DiscogsPutReleasesJson        DiscogsPutReleasesJsonOptions        `command:"releases" description:"Put json for all releases from the release archive in json staging" required:"true"`
	DiscogsPutMastersJson         DiscogsPutMastersJsonOptions         `command:"masters" description:"Put json for all masters from the masters archive in json staging" required:"true"`
}

type MusicBrainzStagingCommands struct {
	MusicBrainzPutArtistsJson         MusicBrainzPutArtistsJsonOptions         `command:"artists" description:"Put json for all artists from the artists archive in json staging" required:"true"`
	MusicBrainzPutArtistsReleasesJson MusicBrainzPutArtistsReleasesJsonOptions `command:"artists-releases" description:"Put artist releases json in json staging for all artists" required:"true"`
	MusicBrainzPutLabelsJson          MusicBrainzPutLabelsJsonOptions          `command:"labels" description:"Put json for all labels from the labels archive in json staging" required:"true"`
	MusicBrainzPutReleasesJson        MusicBrainzPutReleasesJsonOptions        `command:"releases" description:"Put json for all releases from the release archive in json staging" required:"true"`
	MusicBrainzPutMastersJson         MusicBrainzPutMastersJsonOptions         `command:"masters" description:"Put json for all masters from the masters archive in json staging" required:"true"`
}

type IpfsCommands struct {
	PushToIpfs            PushToIpfsOptions            `command:"push-staging-to-ipfs" description:"Push a json staging to ipfs" required:"true"`
	PublishMfsLibraryPath PublishMfsLibraryPathOptions `command:"publish-mfs-path" description:"Publish mfs library path to ipns" required:"true"`
	GetArtistFolderCid    GetArtistFolderCidOptions    `command:"get-artist-folder-cid" description:"Get the cid for an artist folder" required:"true"`
	GetReleaseJsonCid     GetReleaseJsonCidOptions     `command:"get-release-json-cid" description:"Get the cid for a release json document" required:"true"`
}

type UtilsCommands struct {
	ExtractArtistJson             ExtractArtistJsonOptions             `command:"extract-artist-json" description:"Extract as json from artist archive, given disc id" required:"true"`
	ExtractReleaseJson            ExtractReleaseJsonOptions            `command:"extract-release-json" description:"Extract as json from release archive, given disc id" required:"true"`
	MakeArtistsXmlArchive         MakeArtistsXmlArchiveOptions         `command:"make-artists-xml-archive" description:"Make xml archive file containing specified artists"`
	MakeArtistsXmlReleasesArchive MakeArtistsReleasesXmlArchiveOptions `command:"make-artists-releases-xml-archive" description:"Make xml archive file containing specified releases for the artists specified"`
	RunRepoServer                 RepoServerOptions                    `command:"repo-server" description:"Serve repo ipfs json files via api"`
}

// type DiscogsCommands struct {
// }

type IndexingCommands struct {
	DiscogsIndexCmds DiscogsIndexCommands `command:"discogs"`
}

type StagingCommands struct {
	DiscogsStagingCmds     DiscogsStagingCommands     `command:"discogs"`
	MusicBrainzStagingCmds MusicBrainzStagingCommands `command:"musicbrainz"`
}

type Options struct {
	IndexCmds   IndexingCommands `command:"index"`
	StagingCmds StagingCommands  `command:"stage"`
	IpfsCmds    IpfsCommands     `command:"ipfs"`
	UtilsCmds   UtilsCommands    `command:"utils"`
}

func isDefaultIndexer(storeType string) bool {
	if storeType == "sqlite" {
		return true
	} else {
		return false
	}
}

func ctrlcInit() chan struct{} {
	ctrlcChan := make(chan os.Signal, 1)
	signal.Notify(ctrlcChan, os.Interrupt, syscall.SIGTERM)
	done := make(chan struct{}, 1)
	go func() {
		<-ctrlcChan
		close(done)
	}()
	return done
}
