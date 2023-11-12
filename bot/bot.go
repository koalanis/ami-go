package bot

import (
	"fmt"
	"os"
	"strings"

	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	Session *discordgo.Session
}

func DiscordBotInit(discordBotToken string) (*DiscordBot, error) {
	DG, err := discordgo.New("Bot " + discordBotToken)
	// Create a new Discord session using the provided bot token.
	if err != nil {
		fmt.Println("error creating Discord client,", err)
		return nil, err
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	bot := DiscordBot{DG}
	return &bot, nil
}

func ListChannels(bot *DiscordBot, guild string) error {
	err2 := bot.Session.Open()
	if err2 != nil {
		fmt.Println("error opening Discord session,", err2)
		return err2
	}
	listChannels(bot.Session, guild)
	// Cleanly close down the Discord session.
	bot.Session.Close()
	return nil
}

func SendMessage(bot *DiscordBot, msg string, channel string) error {
	err2 := bot.Session.Open()
	if err2 != nil {
		fmt.Println("error opening Discord session,", err2)
		return err2
	}
	bot.Session.ChannelMessageSend(channel, fmt.Sprintf("<@218854888824766475> %s", msg))

	// Cleanly close down the Discord session.
	bot.Session.Close()
	return nil
}

func InteractiveMode(bot *DiscordBot, guild string) error {
	err2 := bot.Session.Open()
	if err2 != nil {
		fmt.Println("error opening Discord session,", err2)
		return err2
	}
	bot.Session.AddHandler(messageCreate)
	listChannels(bot.Session, guild)
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	d := <-sc
	fmt.Printf("\n%s", d)
	fmt.Println()
	// Cleanly close down the Discord session.
	bot.Session.Close()
	return nil
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("start message debug----")
	fmt.Println(m.Content)
	fmt.Println(m.Author)
	fmt.Println(m.Author.ID)
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
		} else if command == "listchannels" {
			guildID := m.GuildID
			channels, err := s.GuildChannels(guildID)
			if err != nil {
				fmt.Println("Error retrieving channels:", err)
				return
			}

			fmt.Println("Channels in the guild:")
			for _, channel := range channels {
				fmt.Printf("ID: %s, Name: %s, Type: %s\n", channel.ID, channel.Name, channel.Type)
			}
		}
	}
}

func listChannels(s *discordgo.Session, guild string) {
	guildID := guild // Replace with your guild ID

	channels, err := s.GuildChannels(guildID)
	if err != nil {
		fmt.Println("Error retrieving channels:", err)
		return
	}

	fmt.Println("Channels in the guild:")
	for _, channel := range channels {
		fmt.Printf("ID: %s, Name: %s, Type: %s\n", channel.ID, channel.Name, channel.Type)
	}
}
