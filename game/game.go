package game

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jordhan-carvalho/belphegorv2/interfaces"
	"github.com/jordhan-carvalho/belphegorv2/sound"
)

var gameTime = 0

func StartListeningToGame(gEventC chan interfaces.GameEvents, vc *discordgo.VoiceConnection, gDone chan bool) {
	fmt.Println("Listening to the game input")
	for {
		select {
		case <-gDone:
			return
		case event := <-gEventC:
			gameTime = event.Map.ClockTime

			checkBountyRunes(vc)
			checkMidRunes(vc)
			checkForStack(vc)
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

func checkForStack(vc *discordgo.VoiceConnection) {
	stackGameTime := 48 // ingame time to stack
	stackDelay := 60    // interval between stack

	if (gameTime-stackGameTime)%stackDelay == 0 {
		go sound.PlaySpecificSound(vc, "stack.dca")
	}
}

func checkBountyRunes(vc *discordgo.VoiceConnection) {
	bountyRunesGameTime := 173 //180
	bountyRunesDelay := 180

	if (gameTime-bountyRunesGameTime)%bountyRunesDelay == 0 {
		go sound.PlaySpecificSound(vc, "bounty-rune.dca")
	}
}

func checkMidRunes(vc *discordgo.VoiceConnection) {
	midRunesGameTime := 110 // 120
	midRunesDelay := 120

	if (gameTime-midRunesGameTime)%midRunesDelay == 0 {
		go sound.PlaySpecificSound(vc, "mid-rune.dca")
	}
}
