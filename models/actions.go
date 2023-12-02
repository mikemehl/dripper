package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mikemehl/dripper/db"
)

type DownloadEpisode string

func DownloadEpisodeAction(d DetailList) tea.Cmd {
	item := d.SelectedItem()
	switch item := item.(type) {
	case *db.Episode:
		return func() tea.Msg { return DownloadEpisode(item.Enclosures[0].URL) }
	}
	return nil
}
