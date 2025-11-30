package main

type GatorError string

func (e GatorError) Error() string{
	return string(e)
}

const NoProgramProvided GatorError = "No program provided. Must start command with 'gator'"
const NotGatorCommandErr GatorError = "command Name does not use 'gator'"
const NoCommandErr GatorError = "no command provided"
const NoCommandExistsErr GatorError = "This gator command does not exists"
const NoArgumentsProvidedErr GatorError = "no arguments provded"
const TooManyArgumentsErr GatorError = "too many arguments for command"