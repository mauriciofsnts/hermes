package discord

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/discord/commands"

	// register slash commands
	_ "github.com/mauriciofsnts/hermes/internal/discord/commands/categories"
)

func Start() error {

	dg, err := discordgo.New("Bot " + config.Hermes.Token)

	if err != nil {
		return err
	}

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = dg.Open()

	if err != nil {
		return err
	}

	err = commands.RegisterModules(dg)

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
