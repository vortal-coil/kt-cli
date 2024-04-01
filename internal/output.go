package internal

import (
	"fmt"
	"log"
)

const (
	ModeLog = iota
	ModePlain
	ModeNoNewline
)

var printMode = ModeLog

func SetPrintMode(mode int) {
	printMode = mode
}

func Print(content string, params ...interface{}) {
	text := fmt.Sprintf(content, params...)

	switch printMode {
	case ModePlain:
		fmt.Println(text)
	case ModeNoNewline:
		fmt.Print(text)
	default:
		log.Println(text)
	}
}
