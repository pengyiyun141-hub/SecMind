package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"secmind/internal/article"
)

//RSS_parser

type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

//Atom_parser

type AtomLink struct {
	Href string `xml:"href,attr"`
}

type Entry struct {
	Title string   `xml:"title"`
	Link  AtomLink `xml:"link"`
}

type AtomFeed struct {
	XMLName xml.Name `xml:"feed"` // 根标签名必须是 feed
	Entries []Entry  `xml:"entry"`
}

//Parse

type Common struct {
	XMLName xml.Name
}

func ParseFeed(reader io.Reader, sourceURLshort string) ([]article.FeedArticle, error) {
	var Xmldata []byte
	Xmldata, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("失败：", err)
		return nil, err
	}

	var common Common
	var Feedarticles []article.FeedArticle

	xml.Unmarshal(Xmldata, &common)

	switch common.XMLName.Local {
	case "rss":
		rssData, err := ParseRSS(Xmldata)
		if err != nil {
			fmt.Println("ParseRss失败：", err)
			return nil, err
		}

		for i, item := range rssData.Channel.Items {
			Feedarticles = append(Feedarticles, article.FeedArticle{Source: sourceURLshort, Id: i + 1, Title: item.Title, Link: item.Link})
		}
		
	case "feed":
		atomData, err := ParseAtom(Xmldata)
		if err != nil {
			fmt.Println("ParseAtom失败：", err)
			return nil, err
		}

		for i, entry := range atomData.Entries {
			Feedarticles = append(Feedarticles, article.FeedArticle{Source: sourceURLshort, Id: i + 1, Title: entry.Title, Link: entry.Link.Href})
		}
		
	default:
		fmt.Println("未知格式")

	}
	return Feedarticles, err
}

func ParseRSS(Xmldata []byte) (RSS, error) {
	var rss RSS
	err := xml.Unmarshal(Xmldata, &rss)
	return rss, err
}

func ParseAtom(Xmldata []byte) (AtomFeed, error) {
	var atom AtomFeed
	err := xml.Unmarshal(Xmldata, &atom)
	return atom, err
}