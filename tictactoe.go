package main

import (
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

func sendBoard(s *dg.Session, m *dg.MessageCreate, board string) {
	rows := strings.Split(board, ",")
	out := "Current board:\n```"
	for i := 0; i < 3; i++ {
		row := rows[i]
		for j := 0; j < 3; j++ {
			pos := row[j]
			out += string(pos)
			if j < 2 {
				out += "|"
			}
		}
		if i < 2 {
			out += "\n-+-+-\n"
		}
	}
	out += "```"
	s.ChannelMessageSend(m.ChannelID, out)
}

func tictactoe(s *dg.Session, m *dg.MessageCreate, args []string) {
	out := "To use tictactoe use one of the following subCommands:\n"
	for k := range newMessageCommands["tictactoe"].subCommands {
		out += k + "\n"
	}
	s.ChannelMessageSend(m.ChannelID, out)
}

func start(s *dg.Session, m *dg.MessageCreate, args []string) {
	if len(args) > 0 {
		opponent := strings.TrimPrefix(strings.TrimSuffix(args[0], ">"), "<@")
		if opponent != m.Author.ID {
			if _, ok := mem["request:"+m.Author.ID+","+m.ChannelID]; !ok {
				mem["request:"+m.Author.ID+","+m.ChannelID] = opponent
				s.ChannelMessageSend(m.ChannelID, "Request sent\nOpponent must accept by doing $tictactoe accept <@"+m.Author.ID+">")
			} else {
				s.ChannelMessageSend(m.ChannelID, "Cannot make multiple requests on the same channel")
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "Cannot send request to yourself")
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please specify a person to start a game with")
	}
}

func cancelRequest(s *dg.Session, m *dg.MessageCreate, args []string) {
	if _, ok := mem["request:"+m.Author.ID+","+m.ChannelID]; ok {
		delete(mem, "request:"+m.Author.ID+","+m.ChannelID)
		s.ChannelMessageSend(m.ChannelID, "Request canceled")
	}
}

func accept(s *dg.Session, m *dg.MessageCreate, args []string) {
	if len(args) > 0 {
		opponent := strings.TrimPrefix(strings.TrimSuffix(args[0], ">"), "<@")
		if v, ok := mem["request:"+opponent+","+m.ChannelID]; ok {
			if v == m.Author.ID {
				delete(mem, "request:"+opponent+","+m.ChannelID)
				emptyBoard := "   ,   ,   "
				mem["game:"+m.Author.ID+","+opponent+","+m.ChannelID] = emptyBoard
				s.ChannelMessageSend(m.ChannelID, "Game started")
				sendBoard(s, m, emptyBoard)
			} else {
				s.ChannelMessageSend(m.ChannelID, "That person has not sent you a request")
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "That person has not sent you a request")
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please specify whose request to accept")
	}
}
