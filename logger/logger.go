package logger

import (
	"fmt"
	"os"
	"log"
	"strings")

var (
	enableFileLogging  = false
	MultipleFiles      bool
	totalLog           []string
	downloadBarStarted bool
)

func SetLogToFile(log bool) {
	enableFileLogging = log
}

// Log function logs the formatted string to the console or to a file depending on if the file logging is enabled.
func Log(pattern string, a ...interface{}) {
	formatted_string := fmt.Sprintf(pattern, a...)
	if !enableFileLogging { // if file logging is not enabled, print the formatted string
		fmt.Print(formatted_string)
	} else {
		totalLog = append(totalLog, formatted_string)
		logToFile()
	}
}

func logToFile() {
	if !enableFileLogging {
		return
	}
	file, err := os.Create("wget-log")
	if err != nil {
		log.Fatalf("\nError creating log file! %s", err)
	}
	defer file.Close()

	_, err = file.WriteString(strings.Join(totalLog, ""))
	if err != nil {
		log.Fatalf("\nError Eriting to log file! %s", err)
	}
}