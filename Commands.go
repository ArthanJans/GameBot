package main

import (
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

var newMessageCommands = map[string]*command{}

func parseCommand(s *dg.Session, m *dg.MessageCreate) {
	if m.Content[0] == '$' {
		message := m.Content[1:]
		words := strings.Fields(message)
		if v, ok := newMessageCommands[words[0]]; ok {
			com := v
			commandLength := 1
			for _, word := range words[1:] {
				if _, ok := com.subCommands[word]; ok {
					com = com.subCommands[word]
					commandLength++
				} else {
					break
				}
			}
			com.mainCommand(s, m, words[commandLength:])
		}
	}
}

func addCommand(call string, function func(s *dg.Session, m *dg.MessageCreate, args []string)) {
	com := command{function, make(map[string]*command)}
	com.add(call)
}

func addSubCommands(call string, function func(s *dg.Session, m *dg.MessageCreate, args []string), parentCommand string) {
	words := strings.Fields(parentCommand)
	if v, ok := newMessageCommands[words[0]]; ok {
		com := v
		commandLength := 1
		for _, word := range words[1:] {
			if _, ok := com.subCommands[word]; ok {
				com = com.subCommands[word]
				commandLength++
			} else {
				break
			}
		}
		com.addSubCommand(call, function)
	}
}

func help(s *dg.Session, m *dg.MessageCreate, args []string) {
	out := "The following commands are available:\n"
	for k := range newMessageCommands {
		out += "$" + k + "\n"
	}
	s.ChannelMessageSend(m.ChannelID, out)
}

func commandSetup() {
	addCommand("gamehelp", help)
	addCommand("tictactoe", tictactoe)
	addSubCommands("start", start, "tictactoe")
	addSubCommands("cancelRequest", cancelRequest, "tictactoe")
}
