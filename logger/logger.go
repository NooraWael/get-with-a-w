package logger

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

// This function will handle the downloading of the specified location passed and will
// create or replace the file content of the filename that was passed using the flag then
// write the progress into the log file with the specified name
//
// @params url - the URL which contents will be downloaded
// @params flagVlaue - the filename that was specified with the flag -B
func DownloadAndLog(url string,logFilee string) {
	fmt.Println("Output will be written to wget-log.")
	logFile, err := os.OpenFile("wget-log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	startTime := time.Now()
	logger.Printf("start at %v", startTime.Format("2006-01-02 15:04:05"))

	resp, err := http.Get(url)
	if err != nil {
		logger.Fatalf("Error downloading link: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Fatalf("Status: %v", resp.Status)
	}
	logger.Println("sending request, awaiting response... status 200 OK")

	size := resp.ContentLength
	sizeMB := float64(size) / (1024 * 1024)
	logger.Printf("content size: %d [~%.2fMB]", size, sizeMB)

	file, err := os.Create("test") // Always save downloaded file as 'test'
	if err != nil {
		logger.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()
	logger.Printf("saving file to: %s", "test")

	// Create a progress bar that writes to discard since we don't need to display it
	bar := progressbar.NewOptions64(size, progressbar.OptionSetWriter(io.Discard))
	multiWriter := io.MultiWriter(file, bar)

	// Copy the response body to the file and update the progress bar
	_, err = io.Copy(multiWriter, resp.Body)
	if err != nil {
		logger.Fatalf("Error writing to file: %v", err)
	}

	finishTime := time.Now()
	logger.Printf("Downloaded [%s]\nfinished at %s", url, finishTime.Format("2006-01-02 15:04:05"))
}