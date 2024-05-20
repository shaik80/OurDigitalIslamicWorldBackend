package logger

import (
	"log"
	"os"
)

// Logger is the global logger for the application
var Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
