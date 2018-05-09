package main

import (
	"fmt"
	"math/rand"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

func inGame(player string, channelID string) bool {
	return getGame(player, channelID) != ""
}

func getGame(player string, channelID string) string {
	for k, v := range mem {
		params := strings.Split(k, ",")
		if params[0] == "game" && params[3] == channelID && (params[1] == player || params[2] == player) {
			return v
		}
	}
	return ""
}

func getOpponent(player string, channelID string) string {
	for k := range mem {
		params := strings.Split(k, ",")
		if params[0] == "game" && params[3] == channelID {
			if params[1] == player {
				return params[2]
			} else if params[2] == player {
				return params[1]
			}
		}
	}
	return ""
}

func sendBoard(s *dg.Session, m *dg.MessageCreate, board string, players []string) {
	rows := strings.Split(board, ",")
	if len(rows) == 4 {
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
		out += "```\n"
		if rows[3] == "0" {
			out += "<@" + players[0] + ">'s turn"
		} else {
			out += "<@" + players[1] + ">'s turn"
		}
		s.ChannelMessageSend(m.ChannelID, out)
	} else {
		fmt.Println("Board format incorrect")
	}
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
		if !inGame(m.Author.ID, m.ChannelID) {
			opponent := strings.TrimPrefix(strings.TrimSuffix(args[0], ">"), "<@")
			if opponent != m.Author.ID {
				if _, ok := mem["request,"+m.Author.ID+","+m.ChannelID]; !ok {
					mem["request,"+m.Author.ID+","+m.ChannelID] = opponent
					s.ChannelMessageSend(m.ChannelID, "Request sent\nOpponent must accept by doing $tictactoe accept <@"+m.Author.ID+">")
				} else {
					s.ChannelMessageSend(m.ChannelID, "Cannot make multiple requests on the same channel")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Cannot send request to yourself")
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "Cannot start a game while you are already in a game")
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please specify a person to start a game with")
	}
}

func cancelRequest(s *dg.Session, m *dg.MessageCreate, args []string) {
	if _, ok := mem["request,"+m.Author.ID+","+m.ChannelID]; ok {
		delete(mem, "request,"+m.Author.ID+","+m.ChannelID)
		s.ChannelMessageSend(m.ChannelID, "Request canceled")
	}
}

func accept(s *dg.Session, m *dg.MessageCreate, args []string) {
	if len(args) > 0 {
		if !inGame(m.Author.ID, m.ChannelID) {
			opponent := strings.TrimPrefix(strings.TrimSuffix(args[0], ">"), "<@")
			if v, ok := mem["request,"+opponent+","+m.ChannelID]; ok {
				if v == m.Author.ID {
					delete(mem, "request,"+opponent+","+m.ChannelID)
					start := rand.Intn(2)
					emptyBoard := "   ,   ,   ," + string(start)
					mem["game,"+m.Author.ID+","+opponent+","+m.ChannelID] = emptyBoard
					s.ChannelMessageSend(m.ChannelID, "Game started")
					sendBoard(s, m, emptyBoard, []string{opponent, m.Author.ID})
				} else {
					s.ChannelMessageSend(m.ChannelID, "That person has not sent you a request")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "That person has not sent you a request")
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "Cannot accept a game while you are already in a game")
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Please specify whose request to accept")
	}
}

func concede(s *dg.Session, m *dg.MessageCreate, args []string) {
	if opponent := getOpponent(m.Author.ID, m.ChannelID); opponent != "" {
		s.ChannelMessageSend(m.ChannelID, "Congratulations <@"+opponent+"> wins!")
		delete(mem, "game,"+m.Author.ID+","+opponent+","+m.ChannelID)
		delete(mem, "game,"+opponent+","+m.Author.ID+","+m.ChannelID)
	}
}
