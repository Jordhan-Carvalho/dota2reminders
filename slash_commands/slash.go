package slash_commands

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jordhan-carvalho/belphegorv2/game"
	"github.com/jordhan-carvalho/belphegorv2/interfaces"
	"github.com/jordhan-carvalho/belphegorv2/sound"
	"github.com/jordhan-carvalho/belphegorv2/utils"
)

var gameDone = make(chan bool)

type SlashCommandsHandler struct {
	VoiceStarted   *bool
	GameEventsChan chan interfaces.GameEvents
}

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	// https://discord.com/developers/docs/interactions/application-commands
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "rita",
			Description: "Volta desgramada! Chora Dieguin!!",
		},
		{
			Name:        "join",
			Description: "Join the voice channel and start listening to the game events",
		},
		{
			Name:        "quit",
			Description: "Quits the voice channel",
		},
		{
			Name:        "time",
			Description: "Gets the game current time",
		},
	}
)

func (h *SlashCommandsHandler) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commandName := i.ApplicationCommandData().Name

	/* ** TIME ** */
	if commandName == "time" {
		gameEvent := <-h.GameEventsChan
		message := utils.SecondsToMinutes(gameEvent.Map.ClockTime)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message,
			},
		})

		time.Sleep(time.Second * 5)
		s.InteractionResponseDelete(i.Interaction)
	}

	/* ** JOIN ** */
	if commandName == "join" {
		g, err := s.State.Guild(i.GuildID)
		if err != nil {
			fmt.Println("Could not find the guild: ", err)
			return
		}

		for _, vs := range g.VoiceStates {
			if vs.UserID == i.Member.User.ID {
				vc, err := s.ChannelVoiceJoin(g.ID, vs.ChannelID, false, true)
				if err != nil {
					fmt.Println("Error joining channel: ", err)
					return
				}

				*h.VoiceStarted = true
				go sound.PlaySpecificSound(vc, "diego.dca")
				go game.StartListeningToGame(h.GameEventsChan, vc, gameDone)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Start listening to game events",
				},
			})

			time.Sleep(time.Second * 5)
			s.InteractionResponseDelete(i.Interaction)
		}

	}

	// ** RITA **
	if commandName == "rita" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Volta desgramada! Chora Dieguin!!",
			},
		})

		g, err := s.State.Guild(i.GuildID)
		if err != nil {
			fmt.Println("Could not find the guild: ", err)
			return
		}

		for _, vs := range g.VoiceStates {
			if vs.UserID == i.Member.User.ID {
				vc, err := s.ChannelVoiceJoin(g.ID, vs.ChannelID, false, true)
				if err != nil {
					fmt.Println("Error joining channel: ", err)
					return
				}

				sound.PlaySpecificSound(vc, "rita.dca")
			}
		}

		s.InteractionResponseDelete(i.Interaction)
	}

	// ** QUIT **
	if commandName == "quit" {
		g, err := s.State.Guild(i.GuildID)
		if err != nil {
			fmt.Println("Could not find the guild: ", err)
			return
		}

		for _, vs := range g.VoiceStates {
			if vs.UserID == i.Member.User.ID {
				vc, err := s.ChannelVoiceJoin(g.ID, vs.ChannelID, false, true)
				if err != nil {
					fmt.Println("Error joining channel: ", err)
					return
				}

				vc.Disconnect()
				log.Println("Left voice channel")
				gameDone <- true
			}
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Bot left the voice channel",
			},
		})

		time.Sleep(time.Second * 5)
		s.InteractionResponseDelete(i.Interaction)
	}
}
