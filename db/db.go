package db

import (
	"bytes"
	"encoding/gob"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/charm/kv"
	"github.com/charmbracelet/log"
	"github.com/gosimple/slug"
	"github.com/mmcdole/gofeed"
)

const dbName = "dripper-data"

type (
	Feed    gofeed.Feed
	Episode gofeed.Item
	SubData struct {
		Feeds    []Feed
		Episodes []*Episode
	}
)

// Implementations for the DetailListItem interface
func (f Feed) FilterValue() string    { return f.Title }
func (f Feed) Name() string           { return f.Title }
func (f Feed) Details() string        { return f.Description }
func (e Episode) FilterValue() string { return e.Title }
func (e Episode) Name() string        { return e.Title }
func (e Episode) Details() string     { return e.Description }

// Helper functions for the wrapper types
func (f Feed) Episodes() []*Episode {
	episodes := make([]*Episode, 0, len(f.Items))
	for i, item := range f.Items {
		episodes[i] = (*Episode)(item)
	}
	return episodes
}

func (s *SubData) LoadFeed(feed gofeed.Feed) {
	s.Feeds = append(s.Feeds, Feed(feed))
	for _, it := range feed.Items {
		s.Episodes = append(s.Episodes, (*Episode)(it))
	}
}

func LoadFeeds() tea.Msg {
	db, _ := kv.OpenWithDefaults(dbName)
	defer db.Close()

	keys, _ := db.Keys()

	var subData SubData
	for _, key := range keys {
		raw_feed, _ := db.Get(key)
		feed, _ := decodeFeed(raw_feed)
		subData.LoadFeed(feed)
	}

	data := SubData(subData)
	log.Debug("Loaded feeds from db")
	return &data
}

func NewFeed(url string) (gofeed.Feed, error) {
	feed, err := newFeedFromURL(url)
	if err != nil {
		return gofeed.Feed{}, err
	}
	db, err := kv.OpenWithDefaults(dbName)
	if err != nil {
		return gofeed.Feed{}, err
	}
	defer db.Close()
	err = addFromFeed(db, feed)
	if err != nil {
		return gofeed.Feed{}, err
	}
	if err := db.Sync(); err != nil {
		return gofeed.Feed{}, err
	}

	return feed, nil
}

func newFeedFromURL(url string) (gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return gofeed.Feed{}, err
	}
	return *feed, nil
}

func addFrromFeeds(feeds []gofeed.Feed) {
	db, _ := kv.OpenWithDefaults(dbName)
	defer db.Close()

	for _, feed := range feeds {
		_ = addFromFeed(db, feed)
	}

	_ = db.Sync()
}

func addFromFeed(db *kv.KV, feed gofeed.Feed) error {
	key := []byte(slug.Make(feed.Title))
	if val, _ := db.Get(key); val != nil {
		return nil // already exists
	}
	val, err := encodeFeed(feed)
	if err != nil {
		return err
	}
	db.Set(key, val)
	return nil
}

func encodeFeed(feed gofeed.Feed) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(feed); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func decodeFeed(data []byte) (gofeed.Feed, error) {
	var feed gofeed.Feed
	dec := gob.NewDecoder(bytes.NewReader(data))
	if err := dec.Decode(&feed); err != nil {
		return gofeed.Feed{}, err
	}
	return feed, nil
}
