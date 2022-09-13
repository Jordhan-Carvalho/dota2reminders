package game

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jordhan-carvalho/belphegorv2/interfaces"
	"github.com/jordhan-carvalho/belphegorv2/sound"
)

var stackGameTime = 48        // ingame time to stack
var stackDelay = 60           // interval between stack
var bountyRunesGameTime = 173 //180
var bountyRunesDelay = 180
var riverRunesGameTime = 110 // 120
var gameTime = 0

func StartListeningToGame(gEventC chan interfaces.GameEvents, vc *discordgo.VoiceConnection, gDone chan bool) {
	fmt.Println("Listening to the game input")
	for {
		select {
		case <-gDone:
			return
		case event := <-gEventC:
			// get the game time
			gameTime = event.Map.ClockTime

			if (gameTime-stackGameTime)%stackDelay == 0 {
				go sound.PlaySpecificSound(vc, "stack.dca")
			}

			if (gameTime-bountyRunesGameTime)%bountyRunesDelay == 0 {
				go sound.PlaySpecificSound(vc, "runa.dca")
			}

			checkForSmokeInShop(vc)
		}
	}
}

// smoke logic, will check if any on invertory if its any it will start a 7 min count
// smoke at every 7 minutes... it starts at 2... max stack is 3... after 7 min you will have max stack
func checkForSmokeInShop(vc *discordgo.VoiceConnection) {
	smokeWarnTime := 410
	smokeDelay := 420
	if (gameTime-smokeWarnTime)%smokeDelay == 0 {
		go sound.PlaySpecificSound(vc, "smoke.dca")
	}
}
