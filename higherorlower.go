package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	dg "github.com/bwmarrin/discordgo"
)

var values = map[int]string{
	1:  "A",
	2:  "2",
	3:  "3",
	4:  "4",
	5:  "5",
	6:  "6",
	7:  "7",
	8:  "8",
	9:  "9",
	10: "10",
	11: "J",
	12: "Q",
	13: "K",
}

var layouts = map[int]string{
	1:  "|          |\n|          |\n|          |\n|    %[1]v    |\n|          |\n|          |\n|          |\n",
	2:  "|    %[1]v    |\n|          |\n|          |\n|          |\n|          |\n|          |\n|    %[1]v    |\n",
	3:  "|    %[1]v    |\n|          |\n|          |\n|    %[1]v    |\n|          |\n|          |\n|    %[1]v    |\n",
	4:  "|  %[1]v  %[1]v  |\n|          |\n|          |\n|          |\n|          |\n|          |\n|  %[1]v  %[1]v  |\n",
	5:  "|  %[1]v  %[1]v  |\n|          |\n|          |\n|    %[1]v    |\n|          |\n|          |\n|  %[1]v  %[1]v  |\n",
	6:  "|  %[1]v  %[1]v  |\n|          |\n|          |\n|  %[1]v  %[1]v  |\n|          |\n|          |\n|  %[1]v  %[1]v  |\n",
	7:  "|  %[1]v  %[1]v  |\n|          |\n|    %[1]v    |\n|  %[1]v  %[1]v  |\n|          |\n|          |\n|  %[1]v  %[1]v  |\n",
	8:  "|  %[1]v  %[1]v  |\n|          |\n|  %[1]v  %[1]v  |\n|          |\n|  %[1]v  %[1]v  |\n|          |\n|  %[1]v  %[1]v  |\n",
	9:  "|  %[1]v  %[1]v  |\n|    %[1]v    |\n|  %[1]v  %[1]v  |\n|          |\n|  %[1]v  %[1]v  |\n|          |\n|  %[1]v  %[1]v  |\n",
	10: "|  %[1]v  %[1]v  |\n|    %[1]v    |\n|  %[1]v  %[1]v  |\n|          |\n|  %[1]v  %[1]v  |\n|    %[1]v    |\n|  %[1]v  %[1]v  |\n",
	11: "|          |\n|          |\n|          |\n|     J    |\n|          |\n|          |\n|          |\n",
	12: "|          |\n|          |\n|          |\n|     Q    |\n|          |\n|          |\n|          |\n",
	13: "|          |\n|          |\n|          |\n|     K    |\n|          |\n|          |\n|          |\n",
}

var suits = map[int]string{
	0: "♥️",
	1: "♣️",
	2: "♦️",
	3: "️♠️",
}

func higherorlower(s *dg.Session, m *dg.MessageCreate, args []string) {
	out := "To use higherorlower use one of the following subCommands:\n"
	for k := range newMessageCommands["higherorlower"].subCommands {
		out += k + "\n"
	}
	s.ChannelMessageSend(m.ChannelID, out)
}

func setupPlayer(playerID string) {
	mem["card,"+playerID] = "1"
	mem["deck,"+playerID] = "1"
	mem["streak,"+playerID] = "0"
	mem["high,"+playerID] = "0"
}

func show(s *dg.Session, m *dg.MessageCreate, args []string) {
	if val, ok := mem["card,"+m.Author.ID]; ok {
		if num, err := strconv.Atoi(val); err != nil {
			fmt.Println(err)
		} else {
			displayCard(s, m, num)
		}
	} else {
		setupPlayer(m.Author.ID)
		displayCard(s, m, 1)
	}
}

func highscore(s *dg.Session, m *dg.MessageCreate, args []string) {
	playerid := ""
	if len(args) == 0 {
		playerid = m.Author.ID
	} else {
		playerid = idFromTag(args[0])
	}
	if val, ok := mem["high,"+playerid]; ok {
		s.ChannelMessageSend(m.ChannelID, val)
	} else {
		setupPlayer(m.Author.ID)
		s.ChannelMessageSend(m.ChannelID, "0")
	}
}

func streak(s *dg.Session, m *dg.MessageCreate, args []string) {
	playerid := ""
	if len(args) == 0 {
		playerid = m.Author.ID
	} else {
		playerid = idFromTag(args[0])
	}
	if val, ok := mem["streak,"+playerid]; ok {
		s.ChannelMessageSend(m.ChannelID, val)
	} else {
		setupPlayer(m.Author.ID)
		s.ChannelMessageSend(m.ChannelID, "0")
	}
}

func playhol(s *dg.Session, m *dg.MessageCreate, args []string) {
	if len(args) > 0 {
		if strings.ToLower(args[0]) == "higher" || strings.ToLower(args[0]) == "lower" {
			var a []int
			lastCard := 0
			if val, ok := mem["card,"+m.Author.ID]; ok {
				if num, err := strconv.Atoi(val); err != nil {
					fmt.Println(err)
				} else {
					lastCard = num
				}
			} else {
				setupPlayer(m.Author.ID)
				lastCard = 1
			}
			if val, ok := mem["deck,"+m.Author.ID]; ok {
				vals := strings.Split(val, ",")
				for i := 1; i < 53; i++ {
					found := false
					for j := range vals {
						if num, err := strconv.Atoi(vals[j]); err != nil {
							fmt.Println(err)
						} else {
							if num == i {
								found = true
							}
						}
					}
					if !found {
						a = append(a, i)
					}
				}
			} else {
				setupPlayer(m.Author.ID)
				for i := 2; i < 53; i++ {
					a = append(a, i)
				}
			}
			if len(a) == 0 {
				a = []int{}
				for i := 1; i < 53; i++ {
					if i != lastCard {
						a = append(a, i)
					}
				}
				mem["deck,"+m.Author.ID] = ""
			}
			card := a[rand.Intn(len(a))]
			displayCard(s, m, card)
			mem["card,"+m.Author.ID] = strconv.Itoa(card)
			mem["deck,"+m.Author.ID] += strconv.Itoa(card)
			if (card-1)%13 > (lastCard-1)%13 && strings.ToLower(args[0]) == "higher" || (card-1)%13 < (lastCard-1)%13 && strings.ToLower(args[0]) == "lower" {
				s.ChannelMessageSend(m.ChannelID, "Congrats you got it right!")
				if num, err := strconv.Atoi(mem["streak,"+m.Author.ID]); err != nil {
					fmt.Println(err)
				} else {
					mem["streak,"+m.Author.ID] = strconv.Itoa(num + 1)
					if high, err := strconv.Atoi(mem["high,"+m.Author.ID]); err != nil {
						fmt.Println(err)
					} else {
						if num+1 > high {
							mem["high,"+m.Author.ID] = strconv.Itoa(num + 1)
						}
					}
				}
			} else if (card-1)%13 > (lastCard-1)%13 && strings.ToLower(args[0]) == "lower" || (card-1)%13 < (lastCard-1)%13 && strings.ToLower(args[0]) == "higher" {
				s.ChannelMessageSend(m.ChannelID, "Unfortunately you got it wrong.")
				mem["streak,"+m.Author.ID] = "0"
			} else {
				s.ChannelMessageSend(m.ChannelID, "It was the same value.")
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "Usage:\n$higherorlower play [higher|lower]")
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, "Usage:\n$higherorlower play [higher|lower]")
	}
}

func displayCard(s *dg.Session, m *dg.MessageCreate, card int) {
	display := "```/¯¯¯¯¯¯¯¯¯¯\\\n"
	suit := card / 13
	value := card % 13
	displayValue, ok := values[value]
	if !ok {
		fmt.Println("Invalid value for display")
	}
	displaySuit, ok := suits[suit]
	if !ok {
		fmt.Println("Invalid suit for display")
	}
	display += "|" + displayValue + "         |\n"
	display += "|" + displaySuit + "        |\n"
	display += layouts[value]
	display += "\\__________/```"
	display = fmt.Sprintf(display, displaySuit)

	s.ChannelMessageSend(m.ChannelID, display)
}
