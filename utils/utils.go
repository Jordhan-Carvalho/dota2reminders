package utils

import "fmt"

func SecondsToMinutes(inSeconds int) string {
	minutes := inSeconds / 60
	seconds := inSeconds % 60
  str := fmt.Sprintf("%02d:%02d\n", minutes, seconds)
	return str
}
