package utils

import (
	tea "github.com/charmbracelet/bubbletea"
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
