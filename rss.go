package main

import "encoding/xml"

type Enclosure struct {
	xml.Name `xml:"enclosure"`
	Length   string `xml:"length,attr"`
	Format   string `xml:"type,attr"`
	Url      string `xml:"url,attr"`
}

type Item struct {
	xml.Name        `xml:"item"`
	Description     string `xml:"description"`
	Enclosure       Enclosure
	Guid            string `xml:"guid"`
	Itunes_duration string `xml:"itunes:duration"`
	Itunes_explicit string `xml:"itunes:explicit"`
	Itunes_image    string `xml:"itunes:image"`
	Link            string `xml:"link"`
	PubDate         string `xml:"pubDate"`
	Title           string `xml:"title"`
}

type Rss struct {
	xml.Name        `xml:"channel"`
	Description     string `xml:"description"`
	Guid            string `xml:"guid"`
	Items           []Item
	Itunes_author   string `xml:"itunes:author"`
	Itunes_category string `xml:"itunes:category"`
	Itunes_complete string `xml:"itunes:complete"`
	Itunes_explicit string `xml:"itunes:explicit"`
	Itunes_image    string `xml:"itunes:image"`
	Itunes_type     string `xml:"itunes:type"`
	Language        string `xml:"language"`
	Link            string `xml:"link"`
	Title           string `xml:"title"`
}

func parseRss(data []byte) (Rss, error) {
	var rss Rss
	err := xml.Unmarshal(data, &rss)
	return rss, err
}
