package cli

import (
	"flag"

	"github.com/bwmarrin/discordgo"
	"github.com/koalanis/ami-go/bot"
)

// returns discordBotToken, cmd, interactive, channel, msg, guild from parsed args
func ParseArgs() (string, string, bool, string, string, string) {
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
	return token, command, runBot, channel, message, guild
}

func HandleCommand(cmd string, msg string, discordSession *discordgo.Session, guildId string, channelId string) {
	if cmd == "list" {
		bot.ListChannels(discordSession, guildId)
	} else if cmd == "msg" {
		bot.SendMessage(discordSession, cmd, channelId)
	}
}
