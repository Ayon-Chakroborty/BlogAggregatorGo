package main

type GatorError string

func (e GatorError) Error() string {
	return string(e)
}

const (
	ErrNoProgramProvided   GatorError = "No program provided. Must start command with 'gator'"
	ErrNotGatorCommand     GatorError = "command Name does not use 'gator'"
	ErrNoCommand           GatorError = "no command provided"
	ErrNoCommandExists     GatorError = "This gator command does not exists"
	ErrNoArgumentsProvided GatorError = "no arguments provded"
	ErrTooManyArguments    GatorError = "too many arguments for command"
	ErrNotEnoughArguments     GatorError = "not enough arguments for command"
	ErrUserNotRegistered   GatorError = "user not registered with gator. Run gator register <name> to register with application"
)
