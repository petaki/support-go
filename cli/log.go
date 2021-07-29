package cli

import (
	"log"
	"os"
)

var (
	// InfoLog log.
	InfoLog = log.New(os.Stdout, Cyan("INFO\t"), log.Ldate|log.Ltime)

	// ErrorLog log.
	ErrorLog = log.New(os.Stderr, Red("ERROR\t"), log.Ldate|log.Ltime|log.Lshortfile)
)
