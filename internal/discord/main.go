package discord

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/discord/events"

	// register slash commands
	_ "github.com/mauriciofsnts/hermes/internal/discord/events/categories"
)

func Start() error {

	dg, err := discordgo.New("Bot " + config.Hermes.Token)

	if err != nil {
		return err
	}

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentDirectMessages)

	err = dg.Open()

	if err != nil {
		return err
	}

	err = events.RegisterModules(dg)

	dg.AddHandler(events.MessageCreate)

	if err != nil {
		return err
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()

	return nil
}
