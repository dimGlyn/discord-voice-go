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

type keyword string

var sounds = []sound{
	sound{make(buffer, 0), dataPath + "EEEEEEEEEEEEEEEEEEEEEE.dca"},
	sound{make(buffer, 0), dataPath + "EIMAI_ENTAKSEI.dca"},
	sound{make(buffer, 0), dataPath + "gamw_tis_katares.dca"},
	sound{make(buffer, 0), dataPath + "re_fyge.dca"},
}

var keywordSound = map[keyword]*sound{
	"ok":       &sounds[0],
	"entaksei": &sounds[1],
	"gamw":     &sounds[2],
	"fyge":     &sounds[3],
}

func loadSounds(key int) error {
	file, err := os.Open(sounds[key].path)
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}
	fmt.Println("load file: ", key, sounds[key].path, file)

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

func (s sound) playSound(session *discordgo.Session, guildID, channelID string) (err error) {
	vc, err := session.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}
	time.Sleep(250 * time.Millisecond)

	vc.Speaking(true)
	for _, buff := range s.b {
		vc.OpusSend <- buff
	}
	vc.Speaking(false)

	time.Sleep(250 * time.Millisecond)
	vc.Disconnect()

	return nil
}

func synonym(word string) keyword {
	switch word {
	case "!ok", "!okey", "!e", "!ee", "!eee":
		return "ok"
	case "!entaksei", "!eimaidaksei", "!daksei":
		return "entaksei"
	case "!fyge", "!fige", "!refigeremalakapodorebro":
		return "fyge"
	case "!gamw", "!katares", "!manoules":
		return "gamw"
	}
	return ""
}
