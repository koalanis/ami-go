package discordUtils

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func DiscordSessionInit(discordBotToken string) (*discordgo.Session, error) {
	session, err := discordgo.New("Bot " + discordBotToken)
	// Create a new Discord session using the provided bot token.
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return nil, err
	}
	return session, nil
}

// Given a discord bot, and a guild (aka a discord server), this function prints a list of channels
func ListChannels(s *discordgo.Session, guild string) {
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

// Given a discord bot, a message and a channelId, this function has the Bot send a message to that channel if it exists
func SendMessage(session *discordgo.Session, msg string, channel string) error {
	err2 := session.Open()
	if err2 != nil {
		fmt.Println("error opening Discord session,", err2)
		return err2
	}
	session.ChannelMessageSend(channel, fmt.Sprintf("<@218854888824766475> %s", msg))

	// Cleanly close down the Discord session.
	session.Close()
	return nil
}
