package sound

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

var soundsBuffers = make(map[string][][]byte)

func PlaySpecificSound(vc *discordgo.VoiceConnection, audioName string) {
	audioBuffers := soundsBuffers[audioName]
	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range audioBuffers {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)
}

func LoadAllSounds() (map[string][][]byte, error) {
	items, _ := ioutil.ReadDir("./sounds_assets/")
	log.Printf("Found %d sounds\n", len(items))

	for _, soundItem := range items {
		soundsBuffers[soundItem.Name()] = make([][]byte, 0)
	}

	for key, value := range soundsBuffers {
		err := loadSound(value, key, soundsBuffers)
		if err != nil {
			fmt.Println("Error loading sound: ", err)
			fmt.Println("Please copy a file.dca to this directory.")
			return nil, err
		}
	}

	return soundsBuffers, nil
}

// loadSound attempts to load an encoded sound file from disk.
func loadSound(sBuffer [][]byte, sName string, soundsBuffers map[string][][]byte) error {
	file, err := os.Open("./sounds_assets/" + sName)

	if err != nil {
		fmt.Println("Something went worng opening audio file:", err)
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		sBuffer = append(sBuffer, InBuf)
		soundsBuffers[sName] = sBuffer
	}
}
