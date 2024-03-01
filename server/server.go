package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/koalanis/ami-go/cli"
	"github.com/koalanis/ami-go/db"
	"github.com/koalanis/ami-go/discordUtils"
)

func printMessageCreateDebugLog(m *discordgo.MessageCreate) {
	fmt.Println("start message debug----")
	fmt.Println(m.Content)
	fmt.Println(m.Author)
	fmt.Println(m.Author.ID)
	fmt.Println(m.ChannelID)
	fmt.Println(m.Mentions)
}

type MessageContext struct {
	channelId string
}

func ServerMessageHandler(server *ServerState) func(*discordgo.Session, *discordgo.MessageCreate) {
	amigoBot := InitBot(server)

	callback := func(s *discordgo.Session, m *discordgo.MessageCreate) {
		printMessageCreateDebugLog(m)
		fmt.Println(s.State.User.Bot)
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

		msgCtx := MessageContext{channelId: m.ChannelID}
		amiCommandContext := createAmiCommandContext(s, m, &msgCtx)
		log.Printf("handling command %t", amiCommandContext.valid)

		if amiCommandContext.valid {
			log.Printf("handling command %s", amiCommandContext.command)
			server.commandInvocationCount += 0
			amigoBot.HandleCommand(amiCommandContext)
		}
	}
	return callback
}

func ServerReactionAddHandler(server *ServerState) func(*discordgo.Session, *discordgo.MessageReactionAdd) {
	callback := func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
		_, err := s.Channel(m.ChannelID)
		if err != nil {
			fmt.Println("Error getting channel data:", err)
			return
		}
		fmt.Printf("%s on message %s\n", m.Emoji.Name, m.MessageID)

	}
	return callback
}

func ServerScheduledActionHandler(server *ServerState) {
	sc := make(chan os.Signal, 0)
	ticker := time.NewTicker((1 * time.Minute))
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		for {
			select {
			case <-ticker.C:
				server.Session.ChannelMessageSend(server.ExecutionContext.Channel, "scheduled message")
			case <-sc:
				ticker.Stop()
				return
			}
		}
	}()
}

func ServerReactionRemoveHandler(server *ServerState) func(*discordgo.Session, *discordgo.MessageReactionRemove) {
	callback := func(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
		channel, err := s.Channel(m.ChannelID)
		if err != nil {
			fmt.Println("Error getting channel data:", err)
			return
		}
		fmt.Println(channel.Name)
		fmt.Println("end message debug----")
		fmt.Printf("%s", m.Emoji.Name)

	}
	return callback
}

// Takes discord bot and creates a server session that listens to messages on the server
func HandleServer(server *ServerState) error {
	err1 := server.Session.Open()
	if err1 != nil {
		fmt.Println("error opening Discord session,", err1)
		return err1
	}

	server.Session.AddHandler(ServerMessageHandler(server))
	server.Session.AddHandler(ServerReactionAddHandler(server))
	server.Session.AddHandler(ServerReactionRemoveHandler(server))

	ServerScheduledActionHandler(server)

	// Cleanly close down the Discord session.
	discordUtils.ListChannels(server.Session, server.ExecutionContext.Guild)
	// Wait here until CTRL-C or other term signal is received.
	defer server.Session.Close()

	waitForExit()

	return nil
}

func waitForExit() {
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 0)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	d := <-sc
	fmt.Printf("\n%s", d)
}

// DiscordBot is a wrapper struct that will contain all metadata for a bot's running session
type ServerState struct {
	Session                *discordgo.Session
	ExecutionContext       cli.AmigoExecutionContext
	commandInvocationCount int
	history                []string
}

func ServerInit(amigo cli.AmigoExecutionContext) error {

	db.InitDb()
	DG, err := discordUtils.DiscordSessionInit(amigo.Token)
	// Create a new Discord session using the provided bot token.
	if err != nil {
		fmt.Println("error creating Discord bot,", err)
		return err
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	server := ServerState{DG, amigo, 0, make([]string, 0)}
	if err != nil {
		return err
	}

	defer fmt.Println("Gracefully shutting down")

	HandleServer(&server)

	return nil
}
