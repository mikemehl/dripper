package pages

import "github.com/mmcdole/gofeed"

func fetchAndParseFeed(url string) (gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	return *feed, err
}
