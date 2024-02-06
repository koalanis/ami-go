package cli

import (
	"flag"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/koalanis/ami-go/bot"
)

type AmigoExecutionContext struct {
	Token      string
	Command    string
	Channel    string
	Message    string
	Guild      string
	ServerMode bool
}

// Parse args and creates the AmigoExecutionContext object
func ParseArgs() AmigoExecutionContext {
	var token string
	var command string
	var channel string
	var message string
	var guild string

	flag.StringVar(&token, "t", "", "Bot Token")
	flag.StringVar(&command, "command", "", "help")
	flag.StringVar(&channel, "channel", "", "s")
	flag.StringVar(&guild, "guild", "", "")
	flag.StringVar(&message, "message", "", "enter message here")

	var runBot bool
	flag.BoolVar(&runBot, "interactive", false, "interactive mode, in which commands are handled by running instance of bot")
	flag.Parse()
	return AmigoExecutionContext{Command: command, Channel: channel, Token: token, Message: message, Guild: guild, ServerMode: runBot}
}

func HandleCommand(cmd string, msg string, discordSession *discordgo.Session, guildId string, channelId string) {
	if cmd == "list" {
		bot.ListChannels(discordSession, guildId)
	} else if cmd == "msg" {
		bot.SendMessage(discordSession, cmd, channelId)
	}
}

func ValidateExecutionContext(amigo AmigoExecutionContext) bool {
	return amigo.Token != "Bot Token"
}

func CliInit(amigo AmigoExecutionContext) {
	discordSession, err := bot.DiscordSessionInit(amigo.Token)
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	HandleCommand(amigo.Command, amigo.Message, discordSession, amigo.Guild, amigo.Channel)
}
