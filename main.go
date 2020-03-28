package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type DiscordPayload struct {
	session *discordgo.Session
	message *discordgo.MessageCreate
}

var (
	Token   string
	scanner *bufio.Scanner
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	discord, err := discordgo.New("Bot " + Token)

	if err != nil {
		fmt.Println("error lol, ", err)
		return
	}

	discord.AddHandler(messageCreate)

	err = discord.Open()
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
	if m.Content == "quit" || m.Content == "exit" {
		fmt.Println("Bye bye")
		return
	}

}

func (dp *DiscordPayload) respond(res string) {
	fmt.Println(res)
	dp.session.ChannelMessageSend(dp.message.ChannelID, res)
}
