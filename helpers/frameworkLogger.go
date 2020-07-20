package helpers

import (
	"log"
	"time"
)

// LogStep calls Output to print to the standard logger with the provided message.
func LogStep(v ...interface{}) {
	log.Println(time.Now().Format("02-Jan-2006"), "---- ", v)
}

// LogError is equivalent to LogStep() followed by a call to os.Exit(1).
func LogError(v ...interface{}) {
	log.Fatalln(time.Now().Format("02-Jan-2006"), "---- ", v)
}
