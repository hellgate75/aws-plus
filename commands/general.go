package commands

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
)

const NilContentText = "--"

var commands = make(map[string]Command)

func CommandNames() []string {
	out := make([]string, 0)
	for k, _ := range commands {
		out = append(out, k)
	}
	return out
}

func CommandByName(name string) (*Command, error) {
	if c, ok := commands[name]; ok {
		return &c, nil
	}
	return nil, errors.New(fmt.Sprintf("Unknown command: %s", name))
}


func ValidateCommand(commandName string) bool {
	if _, ok := commands[commandName]; ok {
		return true
	}
	return false
}

type Arg struct {
	Name	string
	Value	interface{}
}

func (a *Arg) Equals(arg Arg) bool {
	return a.Name == arg.Name
}

func (a *Arg) Int() int {
	return a.Value.(int)
}

func (a *Arg) String() string {
	return a.Value.(string)
}

func (a *Arg) Bool() bool {
	return a.Value.(bool)
}

func (a *Arg) Type() reflect.Type {
	return reflect.TypeOf(a.Value)
}

type ErrorCode	int

const(
	CodeOk					ErrorCode = 0
	CodeBadArgs				ErrorCode = iota + 101
	CodeInsufficientArgs
	CodeTooManyArgs
	CodeParserError
	CodeExecError
	CodeExecFatal
)

type Response struct {
	Code	int			`json:"code,omitempty" yaml:"code,omitempty" xml:"code,omitempty"`
	Error	string		`json:"error,omitempty" yaml:"error,omitempty" xml:"error,omitempty"`
	Content	interface{}	`json:"content,omitempty" yaml:"content,omitempty" xml:"content,omitempty"`
}

type Command struct {
	Name		string
	Sub			[]string
	Args		[]string
	Action		func(*Command, ...Arg) Response
	Parser		func(*Command,[]string, *flag.FlagSet)
	data		map[string]interface{}
}

func (c *Command) Accepts(args ...Arg) ErrorCode {
	if len(c.Args) > len(args) {
		return CodeInsufficientArgs
	}
	for _, arg := range c.Args {
		found := false
		subC:
		for _, xArgs := range args {
			if arg == xArgs.Name {
				found = true
				break subC
			}
		}
		if ! found {
			return CodeBadArgs
		}
	}
	return CodeOk
}

func (c *Command) String() string {
	out := c.Name
	sub := ""
	for _, s := range c.Sub {
		if len(sub) > 0 {
			s += " | "
		}
		sub += s
	}
	if len(sub) > 0 {
		out += " [" + sub + "]"
	}
	for _, a := range c.Args {
		out += " -" + a + "=<" + a + ">"
	}
	return out
}