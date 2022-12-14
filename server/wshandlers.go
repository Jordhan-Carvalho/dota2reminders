package server

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jordhan-carvalho/belphegorv2/game"
	"github.com/jordhan-carvalho/belphegorv2/interfaces"
	"github.com/jordhan-carvalho/belphegorv2/sound"
	"github.com/jordhan-carvalho/belphegorv2/utils"
)

 /* ** DEPRECATED: WS HANDLERS GONE TO SLASH COMMANDS **  */


var gameDone = make(chan bool)
var isRoshanTimerRunning = false

type MessageCreateHandler struct {
	VoiceStarted        *bool
	GameEventsChan      chan interfaces.GameEvents
	GameEventsReceivers *int
	// Vc *discordgo.VoiceConnection
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (h *MessageCreateHandler) Handler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!join" {
		_, g, _ := getChannelAndGuild(s, m)

		// Look for the message sender in that guild's current voice states.
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				// Join the provided voice channel.
				vc, err := s.ChannelVoiceJoin(g.ID, vs.ChannelID, false, true)
				if err != nil {
					fmt.Println("Error joining channel: ", err)
					return
				}

				sound.PlaySpecificSound(vc, "diego.dca")
				*h.VoiceStarted = true

        /* DEPRECATED
				go game.StartListeningToGame(h.GameEventsChan, vc, gameDone) */

				return
			}
		}
	}

	// TODO BROKEN Fix it... it laggs the timer
	if m.Content == "!roshan" {
		// If game is not in progress, do nothing
		event := <-h.GameEventsChan
		fmt.Println("Should be called once")
		// TODO should also check if theres a roshan timer going
		if event.Map.GameState != "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS" || isRoshanTimerRunning {
      fmt.Println("Game not in progress or the time is already running")
			return
		}

		_, g, _ := getChannelAndGuild(s, m)
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				vc, err := s.ChannelVoiceJoin(g.ID, vs.ChannelID, false, true)
				if err != nil {
					fmt.Println("Error joining channel: ", err)
					return
				}

				// TODO: be able to receive a different kill time
				killedTime := event.Map.ClockTime
				go sound.PlaySpecificSound(vc, "roshan-start.dca")
				// TODO best place to increment the receiver?
				*h.GameEventsReceivers++
				isRoshanTimerRunning = true
				go game.StartRoshanAndAegisTimers(h.GameEventsChan, killedTime, vc, &isRoshanTimerRunning, h.GameEventsReceivers)

				return
			}
		}
	}

	if m.Content == "!time" {
		gameEvent := <-h.GameEventsChan
		currentTime := gameEvent.Map.ClockTime
		message := utils.SecondsToMinutes(currentTime)

		_, err := s.ChannelMessageSend(m.ChannelID, message)
		if err != nil {
			fmt.Println(err)
		}
	}

	if m.Content == "!rita" {
		_, g, _ := getChannelAndGuild(s, m)
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				// Join the provided voice channel.
				vc, err := s.ChannelVoiceJoin(g.ID, vs.ChannelID, false, true)
				if err != nil {
					fmt.Println("Error joining channel: ", err)
					return
				}
				go sound.PlaySpecificSound(vc, "rita.dca")
			}
		}
	}

	if m.Content == "!dieguin" {
		_, g, _ := getChannelAndGuild(s, m)
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				// Join the provided voice channel.
				vc, err := s.ChannelVoiceJoin(g.ID, vs.ChannelID, false, true)
				if err != nil {
					fmt.Println("Error joining channel: ", err)
					return
				}
				sound.PlaySpecificSound(vc, "diego-lol.dca")
			}
		}
	}

	if m.Content == "!quit" {
		_, g, _ := getChannelAndGuild(s, m)

		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				// Join the provided voice channel.
				vc, err := s.ChannelVoiceJoin(g.ID, vs.ChannelID, false, true)
				if err != nil {
					fmt.Println("Error joining channel: ", err)
					return
				}

				vc.Disconnect()
				gameDone <- true
				// gameTime = 0
				fmt.Println("Game ended")
			}
		}
	}

}

func getChannelAndGuild(s *discordgo.Session, m *discordgo.MessageCreate) (c *discordgo.Channel, g *discordgo.Guild, err error) {
	// Find the channel that the message came from.
	c, err = s.State.Channel(m.ChannelID)
	if err != nil {
		// Could not find channel.
		return
	}

	// Find the guild for that channel.
	g, err = s.State.Guild(c.GuildID)
	if err != nil {
		// Could not find guild.
		return
	}

	return
}
