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
	"github.com/jordhan-carvalho/belphegorv2/interfaces"
	"github.com/jordhan-carvalho/belphegorv2/server"
	"github.com/jordhan-carvalho/belphegorv2/sound"
)

var token string
var gameEventsChannel = make(chan interfaces.GameEvents)
var gameEventsReceivers = 1
var voiceStarted = false
var port = ":3000"

func init() {
	// This will get the value passed to the program on the flag -t to the token variable
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	if token == "" {
		fmt.Printf("You need to pass the token, please run ./belphegor -t <token value>")
		return
	}

	// Load all sounds in memory
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

	// pass a event and a function to handle the event https://discord.com/developers/docs/topics/gateway#event-names
  messageCreate := &server.MessageCreateHandler{GameEventsChan: gameEventsChannel, VoiceStarted: &voiceStarted, GameEventsReceivers: &gameEventsReceivers}
	discord.AddHandler(messageCreate.Handler)

	// webserver to handler GSI requests
	gameEventsHanlder := &server.GameEventsHandler{GameEventsChan: gameEventsChannel, VoiceStarted: &voiceStarted, GameEventsReceivers: &gameEventsReceivers}
	http.HandleFunc("/", gameEventsHanlder.Handler)

  fmt.Println("Starting http server at port:", port)
  // This way we can clean the discord connection
	go func() {
		if err := http.ListenAndServe(port, nil); err != nil {
			discord.Close()
			log.Fatal(err)
		}
	}()
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()

}
