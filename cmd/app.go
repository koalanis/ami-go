// The app defines the main go process / control flow for the program
// it is responsible for running the program in CLI vs interactive mode
// it also reads the configuration of the bot server before performing assistant actions
package app

import (
	"fmt"
	"os"

	"github.com/koalanis/ami-go/bot"
	"github.com/koalanis/ami-go/cli"
)

func Run() {
	discordBotToken, cmd, interactive, channel, msg, guild := cli.ParseArgs()

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

		cli.HandleCommand(cmd, msg, discordBot, guild, channel)
	}

	os.Exit(0)
}
