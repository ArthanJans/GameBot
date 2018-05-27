package main

import (
	"fmt"
	"strconv"

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
	1:  "|          |\n|          |\n|          |\n|    %v    |\n|          |\n|          |\n|          |\n",
	2:  "|    %v    |\n|          |\n|          |\n|          |\n|          |\n|          |\n|    %v    |\n",
	3:  "|    %v    |\n|          |\n|          |\n|    %v    |\n|          |\n|          |\n|    %v    |\n",
	4:  "|  %v  %v  |\n|          |\n|          |\n|          |\n|          |\n|          |\n|  %v  %v  |\n",
	5:  "|  %v  %v  |\n|          |\n|          |\n|    %v    |\n|          |\n|          |\n|  %v  %v  |\n",
	6:  "|  %v  %v  |\n|          |\n|          |\n|  %v  %v  |\n|          |\n|          |\n|  %v  %v  |\n",
	7:  "|  %v  %v  |\n|          |\n|    %v    |\n|  %v  %v  |\n|          |\n|          |\n|  %v  %v  |\n",
	8:  "|  %v  %v  |\n|          |\n|  %v  %v  |\n|          |\n|  %v  %v  |\n|          |\n|  %v  %v  |\n",
	9:  "|  %v  %v  |\n|    %v    |\n|  %v  %v  |\n|          |\n|  %v  %v  |\n|          |\n|  %v  %v  |\n",
	10: "|  %v  %v  |\n|    %v    |\n|  %v  %v  |\n|          |\n|  %v  %v  |\n|    %v    |\n|  %v  %v  |\n",
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

func show(s *dg.Session, m *dg.MessageCreate, args []string) {
	if val, ok := mem["card,"+m.Author.ID]; ok {
		if num, err := strconv.Atoi(val); err != nil {
			fmt.Println(err)
		} else {
			displayCard(s, m, num)
		}
	} else {
		mem["card,"+m.Author.ID] = "1"
		displayCard(s, m, 1)
	}
}

func playhol(s *dg.Session, m *dg.MessageCreate, args []string) {

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
	display += "\n\\__________/```"
	display = fmt.Sprintf(display, displaySuit)

	s.ChannelMessageSend(m.ChannelID, display)
}
