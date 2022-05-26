package discogstypes

import (
	"encoding/xml"
)

type MinimalArtistXml struct {
	// XMLName  xml.Name `xml:"artist"`
	// InnerXml string   `xml:",innerxml"`
	Id int64 `xml:"id"`
}

type MinimalLabelXml struct {
	// XMLName  xml.Name `xml:"artist"`
	// InnerXml string   `xml:",innerxml"`
	Id int64 `xml:"id"`
}

// type ArtistMember struct {
// 	Id   int64  `xml:"id,attr"`
// 	Name string `xml:"name"`
// }

// func (a *ArtistMember) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
// 	var s string
// 	fmt.Println("DING")
// 	fmt.Println(start)
// 	if err := d.DecodeElement(&s, &start); err != nil {
// 		return err
// 	}
// 	fmt.Println(s)
// 	return nil
// }

type ArtistXmlWithInfo struct {
	Id       int64  `xml:"id"`
	Name     string `xml:"name"`
	RealName string `xml:"realname"`
	// MemberNames []string `xml:"members>name"`
	MemberIds []int64 `xml:"members>id"`
}

type ArtistSearchInfo struct {
	Id       int64
	Name     string
	RealName string
	IsGroup  bool
}

type ReleaseGenreEntry struct {
	Genre string
}

type ReleaseStyleEntry struct {
	Style string
}

type ArtistXML2 struct {
	XMLName  xml.Name `xml:"artist"`
	InnerXml string   `xml:",innerxml"`
	// Id       uint64   `xml:"id"`
}

type ReleaseArtist struct {
	Id int64 `xml:"id"`
	// Name string `xml:"name"`
	// Role string `xml:"role"`
}

type ReleaseLabel struct {
	Id int64 `xml:"id,attr"`
	// Name string `xml:"name"`
	// Role string `xml:"role"`
}

// type Genre struct {
// 	Genre string `xml:"genre"`
// }

type ExtraArtist struct {
	Id int64 `xml:"id"`
	// Name string `xml:"name"`
	// Role string `xml:"role"`
}

type TrackArtist struct {
	Id int64 `xml:"id"`
	// Name string `xml:"name"`
	// Role string `xml:"role"`
}

type ReleaseArtists struct {
	Artist []ReleaseArtist `xml:"artist"`
}

type ReleaseLabels struct {
	Label []ReleaseLabel `xml:"label"`
}

type ReleaseGenres struct {
	Genres []string `xml:"genre"`
}

type ReleaseStyles struct {
	Styles []string `xml:"style"`
}

type ExtraArtists struct {
	Artist []ExtraArtist `xml:"artist"`
}

type TrackArtists struct {
	Artist []TrackArtist `xml:"artist"`
}

type Track struct {
	Position     int          `xml:"position"`
	Title        string       `xml:"title"`
	TrackArtists TrackArtists `xml:"artists"`
}

type TrackList struct {
	Tracks []Track `xml:"track"`
}

// type Title struct {
// 	Title string `xml:"title"`
// }

type MinimalReleaseXml struct {
	// XMLName      xml.Name       `xml:"release"`
	Id int64 `xml:"id,attr"`
	// Id1          uint64         `xml:"id1"`
	// Title string `xml:"title"`
	// Artists      ReleaseArtists `xml:"artists"`
	// Labels       ReleaseLabels  `xml:"labels"`
	// ExtraArtists ExtraArtists   `xml:"extraartists"`
	// TrackList    TrackList      `xml:"tracklist"`
	// InnerXml     string         `xml:",innerxml"`
}

type MinimalMasterXml struct {
	Id int64 `xml:"id,attr"`
}

type ReleaseXmlWithArtists struct {
	Id           int64          `xml:"id,attr"`
	Title        string         `xml:"title"`
	Artists      ReleaseArtists `xml:"artists"`
	ExtraArtists ExtraArtists   `xml:"extraartists"`
	TrackList    TrackList      `xml:"tracklist"`
}

type ReleaseXmlWithLabels struct {
	Id     int64         `xml:"id,attr"`
	Title  string        `xml:"title"`
	Labels ReleaseLabels `xml:"labels"`
}

type ReleaseXmlWithGenres struct {
	Id         int64         `xml:"id,attr"`
	Title      string        `xml:"title"`
	GenresList ReleaseGenres `xml:"genres"`
}

type ReleaseXmlWithStyles struct {
	Id         int64         `xml:"id,attr"`
	Title      string        `xml:"title"`
	StylesList ReleaseStyles `xml:"styles"`
}

type ReleaseXmlWithArtistsAndLabels struct {
	Id           int64          `xml:"id,attr"`
	Title        string         `xml:"title"`
	Artists      ReleaseArtists `xml:"artists"`
	ExtraArtists ExtraArtists   `xml:"extraartists"`
	TrackList    TrackList      `xml:"tracklist"`
	Labels       ReleaseLabels  `xml:"labels"`
}

type XmlReleaseSummary struct {
	XMLName      xml.Name       `xml:"release"`
	Id           int64          `xml:"id"`
	StartPos     int64          `xml:"startpos"`
	EndPos       int64          `xml:"endpos"`
	Title        string         `xml:"title"`
	Artists      ReleaseArtists `xml:"artists"`
	ExtraArtists ExtraArtists   `xml:"extraartists"`
	TrackList    TrackList      `xml:"tracklist"`
	// InnerXml string   `xml:",innerxml"`
}

type MinimalReleaseSummary struct {
	// XMLName      xml.Name       `json:"release"`
	Id       int64 // `json:"id"`
	StartPos int64 // `json:"startPos"`
	EndPos   int64 // `json:"endPos"`
	// Title    string `json:"title"`
	// Artists      ReleaseArtists `json:"artists"`
	// ExtraArtists ExtraArtists   `json:"extraArtists"`
	// TrackList    TrackList      `json:"tracklist"`
	// ReleaseArtists []uint64
	// ReleaseLabels  []uint64
	// ExtraArtists   []uint64
	// TrackArtists   []uint64
}
