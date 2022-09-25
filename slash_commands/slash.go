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

var activeAlerts = &interfaces.ActiveAlerts{}

type SlashCommandsHandler struct {
	VoiceStarted   *bool
	GameEventsChan chan interfaces.GameEvents
}

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer
	alertOptions                   = []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "stack",
			Description: "Boolean option",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "ward",
			Description: "Boolean option",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "bounty-rune",
			Description: "Boolean option",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "smoke",
			Description: "Boolean option",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "mid-rune",
			Description: "Boolean option",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionBoolean,
			Name:        "neutral-items",
			Description: "Boolean option",
			Required:    true,
		},
	}

	// https://discord.com/developers/docs/interactions/application-commands
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "rita",
			Description: "Volta desgramada! Chora Dieguin!!",
		},
		{
			Name:        "join",
			Description: "Join the voice channel and start listening to the game events",
			Options:     alertOptions,
		},
		{
			Name:        "quit",
			Description: "Quits the voice channel",
		},
		{
			Name:        "time",
			Description: "Gets the game current time",
		},
		{
			Name:        "alerts",
			Description: "Show the active alerts",
		},
		{
			Name:        "set-active-alerts",
			Description: "Edit what alerts will be active",
			Options:     alertOptions,
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
		options := i.ApplicationCommandData().Options
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}

		activeAlerts.Stack = optionMap["stack"].BoolValue()
		activeAlerts.Smoke = optionMap["smoke"].BoolValue()
		activeAlerts.Ward = optionMap["ward"].BoolValue()
		activeAlerts.BountyRune = optionMap["bounty-rune"].BoolValue()
		activeAlerts.MidRune = optionMap["mid-rune"].BoolValue()
		activeAlerts.NeutralItems = optionMap["neutral-items"].BoolValue()

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
				go game.StartListeningToGame(h.GameEventsChan, vc, gameDone, *activeAlerts)
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

	/* ** ALERTS OPTIONS ** */
	if commandName == "set-active-alerts" {

		// Access options in the order provided by the user.
		options := i.ApplicationCommandData().Options

		// Or convert the slice into a map
		optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
		for _, opt := range options {
			optionMap[opt.Name] = opt
		}

		margs, msgformat := formatOptionsMessage(optionMap)

		// Pass to the game the new values
		// Quit the running go thread of listening game
		gameDone <- true

		activeAlerts.Stack = optionMap["stack"].BoolValue()
		activeAlerts.Smoke = optionMap["smoke"].BoolValue()
		activeAlerts.Ward = optionMap["ward"].BoolValue()
		activeAlerts.BountyRune = optionMap["bounty-rune"].BoolValue()
		activeAlerts.MidRune = optionMap["mid-rune"].BoolValue()
		activeAlerts.NeutralItems = optionMap["neutral-items"].BoolValue()

		// start a new one passing the options value
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
				go sound.PlaySpecificSound(vc, "diego.dca")
				// increase the receiver
				go game.StartListeningToGame(h.GameEventsChan, vc, gameDone, *activeAlerts)

			}
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// Ignore type for now, they will be discussed in "responses"
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf(
					msgformat,
					margs...,
				),
			},
		})

		time.Sleep(time.Second * 5)
		s.InteractionResponseDelete(i.Interaction)
	}

	if commandName == "alerts" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Those are the alerts %+v \n", *activeAlerts),
			},
		})

		time.Sleep(time.Second * 10)
		s.InteractionResponseDelete(i.Interaction)
	}
}

func formatOptionsMessage(optionsMap map[string]*discordgo.ApplicationCommandInteractionDataOption) (margs []interface{}, msgformat string) {
	// This example stores the provided arguments in an []interface{}
	// which will be used to format the bot's response
	margs = make([]interface{}, 0, len(optionsMap))
	msgformat = "You changed the alerts to receive! " +
		"Take a look at the value(s) you entered:\n"

		// Get the value from the option map.
		// When the option exists, ok = true
	for k := range optionsMap {
		if opt, ok := optionsMap[k]; ok {
			// Option values must be type asserted from interface{}.
			// Discordgo provides utility functions to make this simple.
			margs = append(margs, opt.BoolValue())
			msgformat += "> " + k + ": %v\n"
		}
	}

	return
}
