package logger

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
	"wget/utils"

	"github.com/schollz/progressbar/v3"
)

// This function will handle the downloading of the specified location passed and will
// create or replace the file content of the filename that was passed using the flag then
// write the progress into the log file with the specified name
//
// @params url - the URL which contents will be downloaded
// @params flagVlaue - the filename that was specified with the flag -B
func DownloadAndLog(url string, flagValue string) {
	// identify if the argument provided contains the file name or just the link
	// example:
	// "https://example.com",
	// "www.example.com",
	// "example.com",
	// "ftp://example.com",
	regex := `(?i)^(https?:\/\/)?(www\.)?([a-zA-Z0-9_-]+\.)+[a-zA-Z]{2,}(/.*)?$`
	var outputContainer []string
	// get the file name from the arguement

	// check if the second argument is a filename and not a link
	matched, err := regexp.MatchString(regex, flagValue)
	if err != nil {
		log.Fatalf("Something went wrong: %v", err)
	}

	// check if the filename is valid or not
	if matched {
		log.Fatalf("Invalid filename")
	}
	// open the file in append mode and if it does not exist create it
	logFile, err := os.OpenFile(flagValue, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	// start the timer
	startTime := time.Now()
	outputContainer = append(outputContainer, fmt.Sprintf("start at %v", startTime.Format("006-01-02 15:04:05")))

	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error downloading link: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatalf("Status: %v", err)
	}
	outputContainer = append(outputContainer, "sending request, awaiting response... status 200 OK")

	size := res.ContentLength
	sizeMB := float64(size) / (1024 * 1024)
	outputContainer = append(outputContainer, fmt.Sprintf("content size: %d [~%.2fMB]", size, sizeMB))

	name, err := utils.MakeAName(url)
	if err != nil {
		log.Fatalf("Error making a name for the file: %v", err)
	}

	createdFile, err := os.Create(name)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer createdFile.Close()
	outputContainer = append(outputContainer, fmt.Sprintf("saving file to: %s", name))

	// create the progress bar
	bar := progressbar.DefaultBytesSilent(size, "downloading")

	// Create a multiwriter to write to both file and progress bar
	writer := io.MultiWriter(createdFile, bar)

	// Copy the body to file and progress bar
	_, err = io.Copy(writer, res.Body)
	if err != nil {
		log.Fatalf("Error writing to file: %v", err)
	}

	// get the time it took since the starting time
	finishTime := time.Since(startTime)
	outputContainer = append(outputContainer, fmt.Sprintf("\nDownloaded [%s]\nfinished at %s", url, finishTime))

	// write everything into the log file
	for _, logMessage := range outputContainer {
		if _, err := logFile.WriteString(logMessage + "\n"); err != nil {
			log.Fatalf("Error Writing into log file: %v", err)
		}
	}
}
