package main

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/charmbracelet/charm/kv"
	"github.com/gosimple/slug"
	"github.com/mmcdole/gofeed"
)

const dbName = "dripper-data"

func loadSubData(out chan<- subData) {
	var data subData
	data.feeds = loadFeeds()
	for _, feed := range data.feeds {
		data.episodes = append(data.episodes, feed.Items...)
	}
	out <- data
}

func loadFeeds() []gofeed.Feed {
	db, err := kv.OpenWithDefaults(dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	keys, err := db.Keys()
	if err != nil {
		log.Fatal(err)
	}

	var feeds []gofeed.Feed
	for _, key := range keys {
		raw_feed, err := db.Get(key)
		if err != nil {
			log.Fatal(err)
		}
		feed, err := decodeFeed(raw_feed)
		if err != nil {
			log.Fatal(err)
		}
		feeds = append(feeds, feed)
	}

	return feeds
}

func newFeed(status chan<- error, url string) {
	feed, err := newFeedFromURL(url)
	if err != nil {
		status <- err
		return
	}
	db, err := kv.OpenWithDefaults(dbName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = addFromFeed(db, feed)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Sync(); err != nil {
		log.Fatal(err)
	}
	status <- nil
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
	if val, err := db.Get(key); err != nil {
		return err
	} else if val != nil {
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
