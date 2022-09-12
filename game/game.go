package game

import (
	"fmt"

	"github.com/jordhan-carvalho/belphegorv2/interfaces"
)


func StartGame(gEventC chan interfaces.GameEvents) {
  fmt.Println("Start Game Chamado", gEventC)
  for {
    fmt.Println("Inside loop")
    select {
    case event := <-gEventC:
      fmt.Println(event)
    }
  }
}
