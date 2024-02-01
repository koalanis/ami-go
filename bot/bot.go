// bot package implements reusable functions for creating, and authorizing a discord session
package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// DiscordBot is a wrapper struct that will contain all metadata for a bot's running session
type DiscordBot struct {
	Session                *discordgo.Session
	commandInvocationCount int
	history                []string
}

type MessageContext struct {
	channelId string
}

// function responsible for creating a discord bot
func DiscordBotInit(discordBotToken string) (*DiscordBot, error) {
	DG, err := DiscordSessionInit(discordBotToken)
	// Create a new Discord session using the provided bot token.
	if err != nil {
		fmt.Println("error creating Discord bot,", err)
		return nil, err
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	bot := DiscordBot{DG, 0, make([]string, 0)}
	return &bot, nil
}

func printMessageCreateDebugLog(m *discordgo.MessageCreate) {
	fmt.Println("start message debug----")
	fmt.Println(m.Content)
	fmt.Println(m.Author)
	fmt.Println(m.Author.ID)
	fmt.Println(m.ChannelID)
}

func DiscordBotHelp(bot *DiscordBot, ctx MessageContext) {
	bot.Session.ChannelMessageSend(ctx.channelId, "here is some help")
}

// Takes discord bot and creates a server session that listens to messages on the server
func InteractiveMode(bot *DiscordBot, guild string) error {
	err2 := bot.Session.Open()
	if err2 != nil {
		fmt.Println("error opening Discord session,", err2)
		return err2
	}
	bot.Session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		printMessageCreateDebugLog(m)
		channel, err := s.Channel(m.ChannelID)
		if err != nil {
			fmt.Println("Error getting channel data:", err)
			return
		}
		fmt.Println(channel.Name)
		fmt.Println("end message debug----")

		if m.Author.ID == s.State.User.ID {
			fmt.Println("early return")
			return
		}

		NOT_IMPLEMENTED_YET := func() {
			s.ChannelMessageSend(m.ChannelID, "not implemented yet")
		}

		msgCtx := MessageContext{channelId: m.ChannelID}

		if strings.Contains(m.Content, "!leet") {
			bot.commandInvocationCount += 1
			commands := strings.Split(m.Content, " ")[1:]
			if len(commands) == 0 {
				DiscordBotHelp(bot, msgCtx)
				return
			}
			command := commands[0]
			fmt.Printf("something super cool %s\n", command)
			if command == "help" {
				DiscordBotHelp(bot, msgCtx)
			} else if command == "count" {
				msg := fmt.Sprintf("count = %d", bot.commandInvocationCount)
				fmt.Println(msg)
				s.ChannelMessageSend(m.ChannelID, msg)
				bot.commandInvocationCount -= 1
			} else if command == "todo" {
				s.ChannelMessageSend(m.ChannelID, "here are the todos")
			} else if command == "today" {
				NOT_IMPLEMENTED_YET()
			} else if command == "random" {
				NOT_IMPLEMENTED_YET()
			} else if command == "stats" {
				NOT_IMPLEMENTED_YET()
			}
		}
	})
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
