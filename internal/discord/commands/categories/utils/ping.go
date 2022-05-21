package utils

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mauriciofsnts/hermes/internal/discord/commands"
)

var Ping = &commands.SlashCommand{
	ApplicationCommand: &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Ping the bot",
	},
	Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Pong!",
			},
		})
	},
}
