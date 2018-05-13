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

func setGame(player string, channelID string, game string) {
	for k := range mem {
		params := strings.Split(k, ",")
		if params[0] == "game" && params[3] == channelID && (params[1] == player || params[2] == player) {
			mem[k] = game
		}
	}
	fmt.Println("Error setting game")
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

func sendBoard(s *dg.Session, m *dg.MessageCreate, board string) {
	rows := strings.Split(board, ",")
	if len(rows) == 5 {
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
		if player, err := s.User(rows[3]); err == nil {
			out += player.Username + "'s turn"
		} else {
			fmt.Println("Couldn't find player ID")
			fmt.Println(err)
		}
		s.ChannelMessageSend(m.ChannelID, out)
	} else {
		fmt.Println("Board format incorrect")
	}
}

func checkWin(s *dg.Session, m *dg.MessageCreate, board string) {
	if player, err := s.User(m.Author.ID); err == nil {
		rows := strings.Split(board, ",")
		for i := 0; i < 3; i++ {
			row := rows[i]
			if row == "XXX" || row == "OOO" {
				if opponent := getOpponent(m.Author.ID, m.ChannelID); opponent != "" {
					s.ChannelMessageSend(m.ChannelID, "Congratulations "+player.Username+" wins!")
					delete(mem, "game,"+m.Author.ID+","+opponent+","+m.ChannelID)
					delete(mem, "game,"+opponent+","+m.Author.ID+","+m.ChannelID)
				} else {
					fmt.Println("Game won but no opponent found")
				}
				return
			}
		}
		for i := 0; i < 3; i++ {
			if rows[0][i] == rows[1][i] && rows[1][i] == rows[2][i] {
				if opponent := getOpponent(m.Author.ID, m.ChannelID); opponent != "" {
					s.ChannelMessageSend(m.ChannelID, "Congratulations "+player.Username+" wins!")
					delete(mem, "game,"+m.Author.ID+","+opponent+","+m.ChannelID)
					delete(mem, "game,"+opponent+","+m.Author.ID+","+m.ChannelID)
				} else {
					fmt.Println("Game won but no opponent found")
				}
				return
			}
		}
		if (rows[0][0] == rows[1][1] && rows[1][1] == rows[2][2]) || (rows[0][2] == rows[1][1] && rows[1][1] == rows[2][0]) {
			if opponent := getOpponent(m.Author.ID, m.ChannelID); opponent != "" {
				s.ChannelMessageSend(m.ChannelID, "Congratulations "+player.Username+" wins!")
				delete(mem, "game,"+m.Author.ID+","+opponent+","+m.ChannelID)
				delete(mem, "game,"+opponent+","+m.Author.ID+","+m.ChannelID)
			} else {
				fmt.Println("Game won but no opponent found")
			}
			return
		}
		for x := 0; x < 3; x++ {
			for y := 0; y < 3; y++ {
				if string(rows[0][0]) == " " {
					return
				}
			}
		}
		if opponent := getOpponent(m.Author.ID, m.ChannelID); opponent != "" {
			s.ChannelMessageSend(m.ChannelID, "It's a draw")
			delete(mem, "game,"+m.Author.ID+","+opponent+","+m.ChannelID)
			delete(mem, "game,"+opponent+","+m.Author.ID+","+m.ChannelID)
		} else {
			fmt.Println("Game drawn but no opponent found")
		}
	} else {
		fmt.Println("Couldn't find player ID")
		fmt.Println(err)
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
					emptyBoard := "   ,   ,   ,"
					if start == 0 {
						emptyBoard += m.Author.ID + ",O"
					} else {
						emptyBoard += opponent + ",X"
					}
					mem["game,"+m.Author.ID+","+opponent+","+m.ChannelID] = emptyBoard
					s.ChannelMessageSend(m.ChannelID, "Game started")
					if op, err := s.User(opponent); err == nil {
						if us, err1 := s.User(m.Author.ID); err1 == nil {
							s.ChannelMessageSend(m.ChannelID, op.Username+" is X's and "+us.Username+" is O's")
							sendBoard(s, m, emptyBoard)
						} else {
							fmt.Println("Couldn't find player ID")
							fmt.Println(err1)
						}
					} else {
						fmt.Println("Couldn't find player ID")
						fmt.Println(err)
					}
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
		if player, err := s.User(opponent); err == nil {
			s.ChannelMessageSend(m.ChannelID, "Congratulations "+player.Username+" wins!")
			delete(mem, "game,"+m.Author.ID+","+opponent+","+m.ChannelID)
			delete(mem, "game,"+opponent+","+m.Author.ID+","+m.ChannelID)
		} else {
			fmt.Println("Player conceded but opponent's ID not found")
			fmt.Println(err)
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Cannot concede when not in a game")
	}
}

var rows = map[string]int{
	"top":    0,
	"middle": 1,
	"bottom": 2,
}

var cols = map[string]int{
	"left":   0,
	"middle": 1,
	"right":  2,
}

func play(s *dg.Session, m *dg.MessageCreate, args []string) {
	if len(args) >= 2 {
		if game := getGame(m.Author.ID, m.ChannelID); game != "" {
			if params := strings.Split(game, ","); params[3] == m.Author.ID {
				if row, ok := rows[strings.ToLower(args[0])]; ok {
					if col, ok := cols[strings.ToLower(args[1])]; ok {
						if string(params[row][col]) == " " {
							params[row] = params[row][:col] + params[4] + params[row][col+1:]
							params[3] = getOpponent(m.Author.ID, m.ChannelID)
							if params[4] == "X" {
								params[4] = "O"
							} else {
								params[4] = "X"
							}
							game = strings.Join(params, ",")
							sendBoard(s, m, game)
							setGame(m.Author.ID, m.ChannelID, game)
							checkWin(s, m, game)
						} else {
							s.ChannelMessageSend(m.ChannelID, "That space is not empty")
						}
					} else {
						s.ChannelMessageSend(m.ChannelID, "Usage:\n$tictactoe play [top|middle|bottom] [left|middle|right]")
					}
				} else {
					s.ChannelMessageSend(m.ChannelID, "Usage:\n$tictactoe play [top|middle|bottom] [left|middle|right]")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "It is not your turn")
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "You are not in a game")
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Usage:\n$tictactoe play [top|middle|bottom] [left|middle|right]")
	}
}
