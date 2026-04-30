package main

import (
	"encoding/xml"
	"io"
)

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
	return rss,err
}

