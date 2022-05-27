package utils

import (
	"github.com/bwmarrin/discordgo"
	"github.com/mauriciofsnts/hermes/internal/discord/events"
	"github.com/mauriciofsnts/hermes/internal/utils/reply"
)

func init() {
	events.RegisterCommand(
		&events.SlashCommand{
			ApplicationCommand: &discordgo.ApplicationCommand{
				Name:        "coinflip",
				Description: "Flip a coin",
			},

			Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
				reply.Ok(s, i, &discordgo.MessageEmbed{
					Title: "Coroa!",
					Image: &discordgo.MessageEmbedImage{
						URL: "https://media3.giphy.com/media/afYKFBzlCLx8rmqzFM/giphy.gif?cid=ecf05e47bragd3wljzqt1ib7axzph6khsyqljn96k8n5nx07&rid=giphy.gif&ct=g",
					},
				})
			},
		},
	)
}
