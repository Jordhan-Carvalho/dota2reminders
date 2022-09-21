package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/jordhan-carvalho/belphegorv2/interfaces"
	"github.com/jordhan-carvalho/belphegorv2/server"
	"github.com/jordhan-carvalho/belphegorv2/slash_commands"
	"github.com/jordhan-carvalho/belphegorv2/sound"
)

var (
	token               string
	gameEventsChannel   = make(chan interfaces.GameEvents)
	gameEventsReceivers = 1
	voiceStarted        = false
	port                = ":3000"
	GuildID             string
	RemoveCommands      bool
)

func init() {
	initDotEnv()
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.StringVar(&GuildID, "guild", "", "Test guild ID, Of not passed - bot registers commands globally")
	flag.BoolVar(&RemoveCommands, "rmcmd", true, "Remove all commands after shutdowning or not")
	flag.Parse()
}

func main() {
	if token == "" {
    token = os.Getenv("BOT_TOKEN")
		if token == "" {
			fmt.Printf("You need to pass the token, please run ./belphegor -t <token value> or pass the token through the BOT_TOKEN env var")
			return
		}
	}

	// Load all sounds in memory
	log.Println("Loading all sounds...")
	_, err := sound.LoadAllSounds()
	if err != nil {
		fmt.Println("Error loading the sounds:", err)
	}

	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating the discord session")
		return
	}

	// In this example, we only care about receiving message events.
	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	log.Println("Adding commands...")
	slashCommandsHandler := &slash_commands.SlashCommandsHandler{GameEventsChan: gameEventsChannel, VoiceStarted: &voiceStarted}
	discord.AddHandler(slashCommandsHandler.Handler)
	// it will create a map with the size of the number of commands
	registeredCommands := make([]*discordgo.ApplicationCommand, len(slash_commands.Commands))
	for i, v := range slash_commands.Commands {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	/* // DEPRECATED: In favor of slash commands
	// pass a event and a function to handle the event https://discord.com/developers/docs/topics/gateway#event-names
	messageCreate := &server.MessageCreateHandler{GameEventsChan: gameEventsChannel, VoiceStarted: &voiceStarted, GameEventsReceivers: &gameEventsReceivers}
	discord.AddHandler(messageCreate.Handler) */

	// webserver to handler GSI requests
	gameEventsHanlder := &server.GameEventsHandler{GameEventsChan: gameEventsChannel, VoiceStarted: &voiceStarted, GameEventsReceivers: &gameEventsReceivers}
	http.HandleFunc("/", gameEventsHanlder.Handler)

	log.Println("Starting http server at port", port)
	// This way we can clean the discord connection
	go func() {
		if err := http.ListenAndServe(port, nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Println("Gracefully shutting down.")

	if RemoveCommands {
		log.Println("Removing commands...")
		for _, v := range registeredCommands {
			err := discord.ApplicationCommandDelete(discord.State.User.ID, GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	discord.Close()

}

func initDotEnv() {
	log.Println("Initializing dot env")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
