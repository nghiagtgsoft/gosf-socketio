package logger

import (
	"log"
	"os"

	"github.com/ambelovsky/gosf-socketio/color"
)

func LogDebug(logLine string) {
	if os.Getenv("DEBUG") == "true" {
		log.Println(logLine)
	}
}

func logDebugSocketIo(logLine string) {
	if os.Getenv("DEBUG_SOCKETIO") == "true" {
		log.Println(color.Red + logLine + color.Reset)
	}
}
func logErrorSocketIo(logLine string) {
	log.Println(color.Red + logLine + color.Reset)
}
