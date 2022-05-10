package command

import (
	"strconv"
	"strings"
)

type Command struct {
	Op        string
	Arg       string
	Namespace string
	Text      string
}

// TODO: should this return a pointer to Command (*Command)?
func NewCommand(senderID int, text string) Command {
	var op string
	var arg string

	if strings.HasPrefix(text, "/") {
		fields := strings.Fields(text)
		op = fields[0]
		if len(fields) > 1 {
			arg = fields[1]
		}
	}

	return Command{
		Op:        op,
		Arg:       arg,
		Namespace: strconv.Itoa(senderID),
		Text:      text,
	}
}
