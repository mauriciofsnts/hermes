package commands

import (
	"github.com/bwmarrin/discordgo"
)

type SlashCommand struct {
	*discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var commands = make(map[string]SlashCommand)

func RegisterCommand(command *SlashCommand) {
	commands[command.Name] = *command
}

func RegisterModules(s *discordgo.Session) error {
	applicationCommands := make([]*discordgo.ApplicationCommand, len(commands))

	i := 0
	for _, command := range commands {
		applicationCommands[i] = command.ApplicationCommand
		i++
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commandName := i.ApplicationCommandData().Name

		if command, ok := commands[commandName]; ok {
			command.Handler(s, i)
		}
	})

	_, err := s.ApplicationCommandBulkOverwrite(s.State.User.ID, "", applicationCommands)

	return err
}
