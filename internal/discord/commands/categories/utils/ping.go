package utils

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mauriciofsnts/hermes/internal/discord/commands"
	"github.com/mauriciofsnts/hermes/internal/utils/reply"
)

func init() {
	commands.RegisterCommand(
		&commands.SlashCommand{
			ApplicationCommand: &discordgo.ApplicationCommand{
				Name:        "ping",
				Description: "Ping the bot",
			},
			Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				reply.Ok(s, i, &discordgo.MessageEmbed{
					Description: "Pong!",
				})
			},
		},
	)
}
