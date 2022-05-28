package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (ctx *CommandContext) response(data interface{}) error {
	var embeds []*discordgo.MessageEmbed
	var content string

	switch t := data.(type) {
	case *discordgo.MessageEmbed:
		embeds = append(embeds, t)

	default:
		content = fmt.Sprintf("%v", t)
	}

	return ctx.reply(embeds, content)

}

func (ctx *CommandContext) Error(embed *discordgo.MessageEmbed) {
	embed.Color = 0xe33e32
	ctx.response(embed)
}

func (ctx *CommandContext) Ok(embed *discordgo.MessageEmbed) {
	embed.Color = 0x00ff00
	ctx.response(embed)
}

func (ctx *CommandContext) Text(content string) {
	ctx.response(content)
}

func (ctx *CommandContext) Embed(embed *discordgo.MessageEmbed) {
	ctx.response(embed)
}
