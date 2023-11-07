package main

import (
	"testing"
)

func TestRssParse(t *testing.T) {
	test_rss := []byte(`<?xml version="1.0" encoding="UTF-8"?>
      <rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd">
      <channel>
      <title>Test</title>
      <link>http://example.com</link>
      <description>Test</description>
      <language>en-us</language>
      <itunes:author>Test</itunes:author>
      <itunes:summary>Test</itunes:summary>
      <itunes:explicit>no</itunes:explicit>
      <itunes:name>Test</itunes:name>
      <item>
      <title>Test</title>
      <description>Test</description>
      <enclosure url="http://example.com/test.mp3" length="1" type="audio/mpeg"/>
      <guid>http://example.com/test.mp3</guid>
      <itunes:duration>1:00</itunes:duration>
      <itunes:explicit>no</itunes:explicit>
      <itunes:image href="http://example.com/test.jpg"/>
      <link>http://example.com/test.mp3</link>
      <pubDate>Thu, 01 Jan 1970 00:00:00 +0000</pubDate>
      </item>
      </channel>
      </rss>`)
	parsed, err := parseRss(test_rss)
	if err != nil {
		t.Fatalf("%s", err)
	}
}
