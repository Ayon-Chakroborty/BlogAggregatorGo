package main

type GatorError string

func (e GatorError) Error() string {
	return string(e)
}

const (
	NoProgramProvided      GatorError = "No program provided. Must start command with 'gator'"
	NotGatorCommandErr     GatorError = "command Name does not use 'gator'"
	NoCommandErr           GatorError = "no command provided"
	NoCommandExistsErr     GatorError = "This gator command does not exists"
	NoArgumentsProvidedErr GatorError = "no arguments provded"
	TooManyArgumentsErr    GatorError = "too many arguments for command"
)
