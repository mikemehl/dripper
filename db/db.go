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

func (f Feed) FilterValue() string {
	return f.Title
}

func (e Episode) FilterValue() string {
	return e.Title
}

func (s *SubData) LoadFeed(feed gofeed.Feed) {
	s.Feeds = append(s.Feeds, Feed(feed))
	for _, it := range feed.Items {
		s.Episodes = append(s.Episodes, (*Episode)(it))
	}
}

func LoadFeeds() tea.Msg {
	db, err := kv.OpenWithDefaults(dbName)
	if err != nil {
		log.Fatal("Error opening db: ", "error", err)
	}
	defer db.Close()

	keys, err := db.Keys()
	if err != nil {
		log.Fatal("Error getting keys: ", "error", err)
	} else {
		log.Info("Got keys: ", "keys", keys)
	}

	var subData SubData
	for _, key := range keys {
		raw_feed, err := db.Get(key)
		if err != nil {
			log.Fatal("Error getting key:", "key", key, "error", err)
		}
		feed, err := decodeFeed(raw_feed)
		if err != nil {
			log.Fatal("Error decoding feed:", "error", err)
		}
		log.Debug("Loaded feed: ", "title", feed.Title)
		// log.Debug("Feed has items:", "number of items", len(feed.Items), "items", feed.Items)
		subData.LoadFeed(feed)
	}

	log.Info("Loaded feeds")
	return SubData(subData)
}

func NewFeed(url string) error {
	log.Info("Adding new feed: ", "url", url)
	feed, err := newFeedFromURL(url)
	if err != nil {
		log.Error("Error adding new feed: ", "error", err)
		return err
	}
	db, err := kv.OpenWithDefaults(dbName)
	if err != nil {
		log.Error("Error opening db: ", "error", err)
		return err
	}
	defer db.Close()
	err = addFromFeed(db, feed)
	if err != nil {
		log.Error("Error adding feed to db: ", "error", err)
		return err
	}
	if err := db.Sync(); err != nil {
		log.Info("Error syncing db: ", "error", err)
		return err
	}

	log.Info("Added new feed: ", "title", feed.Title)
	return nil
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
	db, err := kv.OpenWithDefaults(dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for _, feed := range feeds {
		err := addFromFeed(db, feed)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := db.Sync(); err != nil {
		log.Fatal(err)
	}
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
	log.Info("Adding feed: ", "title", feed.Title)
	log.Info("Feed has key: ", "key", key)
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
