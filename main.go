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

func parseToken() string {

	var token string
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
	return token
}

func main() {
	fmt.Println("Hello, Go!")
	discordBotToken := parseToken()
	fmt.Printf("%s\n", discordBotToken)

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + discordBotToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

	fmt.Println("Gracefully shutting down")
	os.Exit(0)
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("start message debug----")
	fmt.Println(m.Content)
	fmt.Println(m.Author)
	fmt.Println(m.ChannelID)
	channel, err := s.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("Error getting channel data:", err)
		return
	}
	fmt.Println(channel.Name)
	fmt.Println("end message debug----\n\n\n")

	if m.Author.ID == s.State.User.ID {
		fmt.Println("early return")
		return
	}

	NOT_IMPLEMENTED_YET := func() {
		s.ChannelMessageSend(m.ChannelID, "not implemented yet")
	}

	if strings.Contains(m.Content, "!leet") {
		commands := strings.Split(m.Content, " ")[1:]
		command := commands[0]
		fmt.Printf("something super cool %s\n", command)
		if command == "help" {
			s.ChannelMessageSend(m.ChannelID, "here is some help")
		} else if command == "today" {
			NOT_IMPLEMENTED_YET()
		} else if command == "random" {
			NOT_IMPLEMENTED_YET()
		} else if command == "stats" {
			NOT_IMPLEMENTED_YET()
		}
	}
}
