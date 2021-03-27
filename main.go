package main

import (
	"fmt"

	"github.com/enescakir/emoji"
	"github.com/manifoldco/promptui"
	"go.uber.org/zap"
)

type ICommand interface {
	Run()

	// CommandName get the name of command to run
	CommandDescription() string
}

type CommandMapFunc map[string]func(*zap.SugaredLogger) ICommand
type CommandMap map[string]ICommand

var commandFuncMap = CommandMapFunc{
	"CREATE_MODELS":     NewGenerateModel,
	"CREATE_REPOSITORY": NewGenerateRepository,
}

func main() {

	logger := NewLogger()
	logger.Infof("%v starting cli ...\n", emoji.Rocket)

	names := []string{}
	commands := CommandMap{}

	for name, command := range commandFuncMap {
		c := command(logger)
		description := c.CommandDescription()
		commandName := fmt.Sprintf("[%s] %s", name, description)
		names = append(names, commandName)
		commands[commandName] = c
	}

	names = append(names, "EXIT")

	prompt := promptui.Select{
		Label: "select command to run",
		Items: names,
	}

	_, result, err := prompt.Run()
	if err != nil {
		logger.Error("prompt failed")
	}

	if result == "EXIT" {
		logger.Infof("application exited %v \n", emoji.PersonBowing)
		return
	}

	commands[result].Run()
}
