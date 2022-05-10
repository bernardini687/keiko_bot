package command

import (
	"log"
)

func Process(cmd Command) string {
	switch cmd.Op {
	case "/get":
		log.Printf("3) cmd: %#v\n", cmd)
		return ReadEntries(cmd)
	case "/stats":
		log.Printf("3) cmd: %#v\n", cmd)
		return Stats(cmd)
	default:
		log.Printf("3) cmd: %#v\n", cmd)
		return AddEntries(cmd)
	}
}
