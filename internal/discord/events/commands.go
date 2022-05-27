package events

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	cfg "github.com/mauriciofsnts/hermes/internal/config"
)

type CommandContext struct {
	Session     *discordgo.Session
	Message     *discordgo.MessageCreate
	Interaction *discordgo.InteractionCreate
	reply       func(embeds []*discordgo.MessageEmbed, content string) error
}

type SlashCommand struct {
	*discordgo.ApplicationCommand
	Alias   []string
	Handler func(ctx *CommandContext)
}

var commands = make(map[string]SlashCommand)

func RegisterCommand(command *SlashCommand) {
	commands[command.Name] = *command

	for _, alias := range command.Alias {
		commands[alias] = *command
	}

}

func RegisterModules(s *discordgo.Session) error {
	var applicationCommands []*discordgo.ApplicationCommand

	for alias, command := range commands {

		if alias != command.Name {
			continue
		}

		applicationCommands = append(applicationCommands, command.ApplicationCommand)
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commandName := i.ApplicationCommandData().Name

		if command, ok := commands[commandName]; ok {
			command.Handler(&CommandContext{
				Session:     s,
				Interaction: i,
				reply: func(embeds []*discordgo.MessageEmbed, content string) error {
					return s.InteractionRespond(
						i.Interaction,
						&discordgo.InteractionResponse{
							Type: discordgo.InteractionResponseChannelMessageWithSource,
							Data: &discordgo.InteractionResponseData{
								Embeds:  embeds,
								Content: content,
							},
						},
					)
				},
			})
		}

	})

	_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", applicationCommands)

	return err
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages from bot himself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// get prefix from config
	prefix := cfg.Hermes.Prefix

	// check if message starts with prefix
	if !strings.HasPrefix(m.Content, prefix) {
		return
	}

	// get command name
	commandName := strings.TrimPrefix(m.Content, prefix)

	cmd, found := commands[commandName]

	if !found {
		return
	}

	cmd.Handler(&CommandContext{
		Session: s,
		Message: m,
		reply: func(embeds []*discordgo.MessageEmbed, content string) error {
			_, err := s.ChannelMessageSendComplex(
				m.ChannelID,
				&discordgo.MessageSend{
					Embeds:  embeds,
					Content: content,
				},
			)

			return err
		},
	})
}
