package utils

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mauriciofsnts/hermes/internal/discord/events"
)

func init() {
	events.RegisterCommand(
		&events.SlashCommand{
			ApplicationCommand: &discordgo.ApplicationCommand{
				Name:        "ping",
				Description: "Ping the bot",
			},
			Alias: []string{"pong", "p"},
			Handler: func(ctx *events.CommandContext) {

				if ctx.Message != nil {
					ctx.Text("PONG MESSAGE ;)")
				} else {
					ctx.Ok(&discordgo.MessageEmbed{
						Title: "Pong!",
					})
				}

			},
		},
	)
}
