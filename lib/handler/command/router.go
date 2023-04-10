package command

import (
	"log"
)

func Process(cmd Command) string {
	log.Printf("3) cmd: %#v\n", cmd)

	switch cmd.Op {
	case "/dues":
		return Dues(cmd)
	case "/get":
		return ReadEntries(cmd)
	case "/stats":
		return Stats(cmd)
	default:
		return AddEntries(cmd)
	}
}
