package db

import (
	"bytes"
	"encoding/gob"
	"slices"
	"strings"

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
	slices.SortFunc(s.Episodes, episodePtrSort)
}

func LoadFeeds() tea.Msg {
	db, _ := kv.OpenWithDefaults(dbName)
	defer db.Close()

	_ = db.Sync()

	keys, _ := db.Keys()

	var subData SubData
	for _, key := range keys {
		raw_feed, err := db.Get(key)
		if err != nil {
			log.Error(err)
			continue
		}
		feed, err := decodeFeed(raw_feed)
		if err != nil {
			log.Error(err)
			continue
		}
		log.Debug("Loaded feed from db: ", "title", feed.Title, "num_items", len(feed.Items))
		subData.LoadFeed(feed)
	}

	data := SubData(subData)
	slices.SortFunc(data.Feeds, feedSort)
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

func UpdateFeeds() tea.Msg {
	db, _ := kv.OpenWithDefaults(dbName)
	defer db.Close()

	_ = db.Sync()

	keys, _ := db.Keys()
	subData := SubData{}
	for _, key := range keys {
		feed_stored, err := db.Get(key)
		if err != nil {
			log.Error("Unable to get feed from db", "key", key, "error", err)
			continue
		}
		feed, err := decodeFeed(feed_stored)
		if err != nil {
			log.Error("Error decoding feed", "key", key, "error", err)
			continue
		}
		feed, err = newFeedFromURL(feed.FeedLink)
		if err != nil {
			log.Error("Error fetching feed", "link", feed.Link, "error", err)
			continue
		}
		dbFeed, err := encodeFeed(feed)
		if err != nil {
			log.Error("Error encoding feed", "feed", feed, "error", err)
			continue
		}
		db.Set(key, dbFeed)
		subData.Feeds = append(subData.Feeds, Feed(feed))
		for _, episode := range feed.Items {
			subData.Episodes = append(subData.Episodes, (*Episode)(episode))
		}
	}
	slices.SortFunc(subData.Feeds, feedSort)
	return &subData
}

func newFeedFromURL(url string) (gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return gofeed.Feed{}, err
	}
	return *feed, nil
}

func addFromFeeds(feeds []gofeed.Feed) {
	db, _ := kv.OpenWithDefaults(dbName)
	defer db.Close()

	for _, feed := range feeds {
		_ = addFromFeed(db, feed)
	}

	_ = db.Sync()
}

func addFromFeed(db *kv.KV, feed gofeed.Feed) error {
	key := []byte(slug.Make(feed.Title))
	if val, err := db.Get(key); val != nil {
		return nil // already exists
	} else if err != nil {
		return err
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

func feedSort(a, b Feed) int {
	return strings.Compare(a.Title, b.Title)
}

func episodeSort(a, b Episode) int {
	if a.PublishedParsed.Before(*b.PublishedParsed) {
		return 1
	} else if a.PublishedParsed.After(*b.PublishedParsed) {
		return -1
	}
	return 0
}

func episodePtrSort(a, b *Episode) int {
	if a.PublishedParsed.Before(*b.PublishedParsed) {
		return 1
	} else if a.PublishedParsed.After(*b.PublishedParsed) {
		return -1
	}
	return 0
}
