package models

import (
	"mime"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/mikemehl/dripper/db"
	utils "github.com/mikemehl/dripper/utils"
)

type DownloadEpisode struct {
	Url      string
	Filename string
}

func DownloadEpisodeAction(d DetailList) tea.Cmd {
	item := d.SelectedItem()
	switch item := item.(type) {
	case *db.Episode:
		return func() tea.Msg {
			extensions, err := mime.ExtensionsByType(item.Enclosures[0].Type)
			if err != nil {
				log.Error("Error getting extension", "error", err)
			}
			return DownloadEpisode{
				Url:      item.Enclosures[0].URL,
				Filename: utils.SlugifyEpisodeTitle(item) + extensions[0],
			}
		}
	}
	return nil
}
