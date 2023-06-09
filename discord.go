package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func Discord() {
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_KEY"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		fmt.Println(m.Author.Username + " | " + m.Content)
		if m.Author.ID != s.State.User.ID && m.Content != "" {

			res, err := AIResponse(m.Content)
			if err != nil {
				panic(err)
			}

			if ch, err := s.State.Channel(m.ChannelID); err != nil || !ch.IsThread() {
				_, _ = s.ChannelMessageSend(m.ChannelID, res)
			} else {
				_, _ = s.ChannelMessageSendReply(m.ChannelID, res, m.Reference())
			}
		}

	})

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

// refferec: https://github.com/bwmarrin/discordgo
