package utils

import (
	"io"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/gosimple/slug"
	"github.com/mikemehl/dripper/db"
)

func ScaleDimensions(msg tea.WindowSizeMsg, bottom int, wtop int, htop int) tea.WindowSizeMsg {
	msg.Width = msg.Width / 10 * wtop
	msg.Height = msg.Height / 10 * htop
	return msg
}

func SlugifyEpisodeTitle(episode *db.Episode) string {
	return slug.Make(episode.Title)
}

func GetAudio(url string) (io.ReadCloser, error) {
	log.Debug("Downloading", "url", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Error("Download failed", "error", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Error("Download failed", "status", resp.StatusCode)
		return nil, err
	}

	log.Debug("Download complete")
	return resp.Body, nil
}
