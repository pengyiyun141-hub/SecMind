package main

import (
	"encoding/xml"
	"io"
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

func ParseRSS(reader io.Reader) (RSS, error){
	var rss RSS
	err := xml.NewDecoder(reader).Decode(&rss)
	if rss.XMLName!= nil {
		rss,err = ParseAtom(reader)
	}
	return rss,err
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


func ParseAtom(reader io.Reader) (AtomFeed, error){
	var atom AtomFeed
	err := xml.NewDecoder(reader).Decode(&atom)
	return atom,err
}
