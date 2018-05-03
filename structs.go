package main

import dg "github.com/bwmarrin/discordgo"

type command struct {
	mainCommand func(s *dg.Session, m *dg.MessageCreate, args []string)
	subCommands map[string]*command
}

func (com command) add(call string) {
	newMessageCommands[call] = &com
}
func (com *command) addSubCommand(call string, subCommand func(s *dg.Session, m *dg.MessageCreate, args []string)) {
	com.subCommands[call] = &command{subCommand, make(map[string]*command)}
}
