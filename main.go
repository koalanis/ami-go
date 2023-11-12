package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/koalanis/ami-go/bot"
)

func parseCli() (string, string, bool, string, string, string) {
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

func main() {
	discordBotToken, cmd, interactive, channel, msg, guild := parseCli()

	if discordBotToken == "Bot Token" {
		os.Exit(1)
		return
	}

	if interactive {
		discordBot, err := bot.DiscordBotInit(discordBotToken)
		if err != nil {
			fmt.Println("error opening connection,", err)
			return
		}
		bot.InteractiveMode(discordBot, guild)
		fmt.Println("Gracefully shutting down")
		os.Exit(0)
		return
	} else {
		discordBot, err := bot.DiscordBotInit(discordBotToken)
		if err != nil {
			fmt.Println("error opening connection,", err)
			return
		}

		if cmd == "list" {
			bot.ListChannels(discordBot, guild)
		} else if cmd == "msg" {
			bot.SendMessage(discordBot, msg, channel)
		}
	}

	os.Exit(0)
}
