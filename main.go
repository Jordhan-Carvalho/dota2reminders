package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/jordhan-carvalho/belphegorv2/server"
	"github.com/jordhan-carvalho/belphegorv2/sound"
	"github.com/jordhan-carvalho/belphegorv2/interfaces"
)

var token string

var stackTime = 48        // ingame time to stack
var stackDelay = 60       // interval between stack
var bountyRunesTime = 173 //180
var bountyRunesDelay = 180
var riverRunesTime = 110 // 120
var gameTime = 0         // SHOULD BE A CHANNEL

var gameEventsChannel = make(chan interfaces.GameEvents)
var voiceStarted = false


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
	soundsBuffers, err := sound.LoadAllSounds()
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
  // go carai(gameEventsChannel)

	// pass a event and a function to handle the event https://discord.com/developers/docs/topics/gateway#event-names
  messageCreate := &server.MessageCreateHandler{SoundsBuffers: soundsBuffers, VoiceStarted: &voiceStarted}
	discord.AddHandler(messageCreate.Handler)

	// webserver to handler GSI requests
  gameEventsHanlder := &server.GameEventsHandler{GameEventsChan: gameEventsChannel, VoiceStarted: &voiceStarted}
	http.HandleFunc("/", gameEventsHanlder.Handler)

	fmt.Printf("Starting server at port 3000\n")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		discord.Close()
		log.Fatal(err)
	}
}



func carai(c chan interfaces.GameEvents) {
  for {
    select {
    case event := <-c:
      fmt.Println(event)
    }
  }
}

