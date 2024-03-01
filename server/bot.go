package server

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/koalanis/ami-go/db"
)

type AmiBot struct {
	ServerState *ServerState
	commands    AmiCommandSpec
}

func isAmiRequest(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	isBotMentioned := false
	for _, v := range m.Mentions {
		if v.ID == s.State.User.ID {
			isBotMentioned = true
		}
	}

	return isBotMentioned || strings.Contains(m.Content, "!ami")
}

type AmiCommandContext struct {
	valid   bool
	msgCtx  *MessageContext
	command string
	args    []string
}

type AmiCommandHandler func(server *ServerState, commandContext AmiCommandContext)
type AmiCommandSpec map[string]AmiCommandHandler

func createAmiCommandContext(s *discordgo.Session, m *discordgo.MessageCreate, msgCtx *MessageContext) AmiCommandContext {
	isRequest := isAmiRequest(s, m)
	tokens := strings.Split(m.Content, " ")[1:]
	if !isRequest || len(tokens) == 0 {
		return AmiCommandContext{valid: false}
	}

	cmd := tokens[0]
	args := tokens[1:]
	return AmiCommandContext{msgCtx: msgCtx, valid: isRequest, command: cmd, args: args}
}

func handleTodo(server *ServerState, commandContext AmiCommandContext) {
	if len(commandContext.args) == 0 {
		return
	}
	if commandContext.args[0] == "list" {
		server.Session.ChannelMessageSend(commandContext.msgCtx.channelId, fmt.Sprintf("here are the todos: %s", db.GetTodos()))
	} else if commandContext.args[0] == "add" {
		db.AddTodo(strings.Join(commandContext.args[0:], " "))
	} else if commandContext.args[0] == "do" {
		if len(commandContext.args[1]) > 0 {
			log.Printf("trying to delete %s", commandContext.args[3])
			db.DoTodo(commandContext.args[3])
		}
	} else if commandContext.args[2] == "clear" {
		db.ClearTodos()
	}
}

func createCommandSpec() AmiCommandSpec {
	spec := make(AmiCommandSpec)

	NOT_IMPLEMENTED_YET := func(commandName string) func(server *ServerState, commandContext AmiCommandContext) {
		return func(ss *ServerState, cc AmiCommandContext) {
			ss.Session.ChannelMessageSend(cc.msgCtx.channelId, fmt.Sprintf("`%s` is not implemented yet", commandName))
		}
	}

	spec["cat"] = func(server *ServerState, commandContext AmiCommandContext) {
		server.Session.ChannelMessageSend(commandContext.msgCtx.channelId, "https://cataas.com/cat")
	}
	spec["todo"] = handleTodo
	spec["today"] = NOT_IMPLEMENTED_YET("today")

	return spec
}

func InitBot(s *ServerState) AmiBot {
	return AmiBot{commands: createCommandSpec(), ServerState: s}
}

func (ami *AmiBot) Help(msgCtx *MessageContext) {
	availableCommands := []string{}
	for key := range ami.commands {
		availableCommands = append(availableCommands, fmt.Sprintf("`%s`", key))
	}
	helpText := fmt.Sprintf("I support the following commands:\n%s", strings.Join(availableCommands, "\n"))
	ami.ServerState.Session.ChannelMessageSend(msgCtx.channelId, helpText)
}

func (ami *AmiBot) HandleCommand(command AmiCommandContext) {

	if command.valid {
		if handler, inMap := ami.commands[command.command]; inMap {
			handler(ami.ServerState, command)
		} else if command.command == "help" {
			ami.Help(command.msgCtx)
		}
	}
}
