package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

type buffer [][]byte

type sound struct {
	b    buffer
	path string
}

var dataPath = "data/"

var sounds = []sound{
	sound{make(buffer, 0), dataPath + "EEEEEEEEEEEEEEEEEEEEEE.dca"},
	sound{make(buffer, 0), dataPath + "EIMAI_ENTAKSEI.dca"},
	sound{make(buffer, 0), dataPath + "gamw_tis_katares.dca"},
}

var keys = map[string]int{
	"ok":       0,
	"entaksei": 1,
	"gamw":     2,
}

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	discord, _ := discordgo.New("Bot " + Token)

	for _, val := range keys {
		err := loadSounds(val)

		if err != nil {
			fmt.Println("error lol, ", err)
			return
		}
	}

	discord.AddHandler(messageCreate)

	err := discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!ok" || m.Content == "!entaksei" || m.Content == "!gamw" {
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			return
		}
		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			return
		}
		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				key := strings.Replace(m.Content, "!", "", 0)
				err = sounds[keys[key]].b.playSound(s, g.ID, vs.ChannelID)
				if err != nil {
					fmt.Println("Error playing sound:", err)
				}

				return
			}
		}
	}
}

func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "Airhorn is ready! Type !airhorn while in a voice channel to play a sound.")
			return
		}
	}
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
