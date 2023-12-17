package player

import (
	"github.com/charmbracelet/log"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	models "github.com/mikemehl/dripper/models"
	"github.com/mikemehl/dripper/utils"
)

func Init() (chan models.DownloadEpisode, error) {
	if err := speaker.Init(44100, 44100/20); err != nil {
		return nil, err
	}
	inbox := make(chan models.DownloadEpisode)
	go processDownloadEpisodes(inbox)
	return inbox, nil
}

func processDownloadEpisodes(inbox chan models.DownloadEpisode) {
	control := make(chan bool)

	select {
	case req := <-inbox:
		audio, err := utils.GetAudio(req.Url)
		if err != nil {
			log.Error("Error downloading audio", "error", err)
			break
		}

		streamer, _, err := mp3.Decode(audio)
		if err != nil {
			log.Error("Error decoding audio", "error", err)
			break
		}

		go playAudio(control, streamer)

	default:
		log.Debug("No play request")
	}
}

func playAudio(control chan bool, streamer beep.Streamer) {
	finished := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		finished <- true
	})))

	for {
		select {
		case act := <-control:
			if act {
				break
			}
		case _ = <-finished:
			break
		default:
		}
	}
}
