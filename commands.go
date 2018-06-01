package main

import (
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

var newMessageCommands = map[string]*command{}

var aliases = map[string]string{}

func parseCommand(s *dg.Session, m *dg.MessageCreate) {
	if m.Content != "" && m.Content[0] == '$' {
		message := m.Content[1:]
		words := strings.Fields(message)
		v, ok := newMessageCommands[words[0]]
		if !ok {
			if alias, ok1 := aliases[words[0]]; ok1 {
				if val, ok2 := newMessageCommands[alias]; ok2 {
					v = val
					ok = true
				}
			}
		}
		if ok {
			com := v
			commandLength := 1
			for _, word := range words[1:] {
				if val, ok := com.subCommands[word]; ok {
					com = val
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

func addSubCommand(call string, function func(s *dg.Session, m *dg.MessageCreate, args []string), parentCommand string) {
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

func addAlias(alias string, fullName string) {
	aliases[alias] = fullName
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
	addAlias("ttt", "tictactoe")
	addSubCommand("start", start, "tictactoe")
	addSubCommand("cancelRequest", cancelRequest, "tictactoe")
	addSubCommand("accept", accept, "tictactoe")
	addSubCommand("concede", concede, "tictactoe")
	addSubCommand("play", playttt, "tictactoe")
	addCommand("higherorlower", higherorlower)
	addAlias("hol", "higherorlower")
	addSubCommand("show", show, "higherorlower")
	addSubCommand("play", playhol, "higherorlower")
	addSubCommand("highscore", highscore, "higherorlower")
	addSubCommand("streak", streak, "higherorlower")
}
