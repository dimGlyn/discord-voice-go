package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

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
				key := strings.Replace(m.Content, "!", "", 1)
				fmt.Println(keys, key, keys[key])
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
