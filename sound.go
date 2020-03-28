package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

type buffer [][]byte

type sound struct {
	b    buffer
	path string
}

func loadSounds(key int) error {
	file, err := os.Open(sounds[key].path)
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		err = binary.Read(file, binary.LittleEndian, &opuslen)

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

		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		sounds[key].b = append(sounds[key].b, InBuf)
	}
}

func (b buffer) playSound(s *discordgo.Session, guildID, channelID string) (err error) {
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}
	time.Sleep(250 * time.Millisecond)

	vc.Speaking(true)
	for _, buff := range b {
		vc.OpusSend <- buff
	}
	vc.Speaking(false)

	time.Sleep(250 * time.Millisecond)
	vc.Disconnect()

	return nil
}
