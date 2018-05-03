package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	dg "github.com/bwmarrin/discordgo"
)

var helpcommand = "$gamehelp"

func main() {
	runtime.GOMAXPROCS(2)
	readJSON(&cfg, "config.json")
	readJSON(&mem, "memory.json")
	commandSetup()
	defer writeJSON(&mem, "memory.json")

	session, err := dg.New(cfg["BotID"])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer session.Close()

	if err = session.Open(); err != nil {
		fmt.Println(err)
		return
	}

	session.AddHandler(newMessage)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV, syscall.SIGHUP)
	fmt.Println("Setup Complete")
	<-sc
}

func newMessage(s *dg.Session, m *dg.MessageCreate) {
	if !m.Author.Bot {
		parseCommand(s, m)
	}
}
