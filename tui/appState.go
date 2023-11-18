package tui

import (
	"github.com/mikemehl/dripper/models"
)

type AppState int

const (
	AppStateInvalid  = -1
	AppStatePodcasts = iota
	AppStateEpisodes
	AppStatePodEpisodes
)

type AppStateError struct {
	msg string
}

func (e *AppStateError) Error() string {
	return e.msg
}

func (a AppState) Update(menu models.Menu) (AppState, error) {
	switch selected := menu.Active(); selected {
	case AppStatePodcasts:
		return AppStatePodcasts, nil
	case AppStateEpisodes:
		return AppStateEpisodes, nil
	case AppStatePodEpisodes:
		return AppStatePodEpisodes, nil
	default:
		return AppStateInvalid, &AppStateError{msg: "Invalid state"}
	}
}

func (a AppState) Render() (string, error) {
	switch a {
	case AppStatePodcasts:
		return "Podcasts", nil
	case AppStateEpisodes:
		return "Episodes", nil
	case AppStatePodEpisodes:
		return "Podcast Episodes", nil
	default:
		return "Error", &AppStateError{msg: "Invalid state"}
	}
}
