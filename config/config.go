package config

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"wget/utils"

	"github.com/schollz/progressbar/v3"
)

//configuration settings and command line flags
//here we call this when there are flags

// bery bery complicated but fun
//
// Handles downloading the file from the specified URL with the given flags and makes them work together
// depending on the flag the log output will change from stdout to a log file. It will make a request
// to the given URL and checks the status of the request and writes to the logWriter
//
// @parameter url - the URL which the file will be downloaded from
// @parameter flags - a map of the entire flags and their values that were passed when running the program
func HandleDownloadWithFlags(url string, flags map[string]string) {
	var err error

	// ----------- -B flag -------------
	logToFile := false // check the existance of the -B flag

	// ----------- -O flag -------------
	changeFileName := false // if the -O flag is passed change the file name that will be created
	var fileName string     // store the value of the file name

	var logger *log.Logger
	var logWriter io.Writer // variable that will handle where the log will be printed to

	for key, value := range flags {
		switch key {
		case "B":
			logToFile = true
			fmt.Println("Output will be written to wget-log if logToFile is true, else to stdout.")
		case "O":
			changeFileName = true
			fileName = value
			fmt.Println(fileName)
		}
	}

	if logToFile {
		// if the log flag is passed open the log file and send the contents to it
		logFile, err := os.OpenFile("wget-log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}

		defer logFile.Close()
		logWriter = logFile
	} else {
		// if the log to file flag is not passed make print the logs into the terminal
		logWriter = os.Stdout
	}

	logger = log.New(logWriter, "", log.LstdFlags)

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

	if !changeFileName {
		// if the change file name flag was not specified make a random name
		fileName, err = utils.MakeAName(url)
		if err != nil {
			logger.Fatalf("Error creating file: %v", err)
		}
	}

	file, err := os.Create(fileName) // Always save downloaded file as 'test'
	if err != nil {
		logger.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()
	logger.Printf("saving file to: ./%s", fileName)

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

func ParseFlags() (map[string]string, bool, bool, string) {
	// Define all possible flags
	outputFileName := flag.String("O", "", "Specify the output file name (optional)")
	downloadPath := flag.String("P", "", "Specify the path to save the file")
	logToFile := flag.Bool("B", false, "Specifiy the filename to write the log into")
	rateLimit := flag.String("rate-limit", "", "Specify the maximum download rate (e.g., '500k', '2M')")
	help := flag.Bool("help", false, "Display help information")
	web := flag.Bool("web", false, "Start the web server interface")

	// Parse the command line arguments
	flag.Parse()

	// Check if help was requested
	if *help {
		utils.DisplayHelp()
		os.Exit(0)
	}

	// Map to store flags that were actually set
	flagsUsed := make(map[string]string)
	anyFlagUsed := false

	// captures the non flag arguements
	url := flag.Arg(0)

	// Check each flag and add to the map if it was set
	if *outputFileName != "" {
		flagsUsed["O"] = *outputFileName
		anyFlagUsed = true
	}

	if *downloadPath != "" {
		flagsUsed["P"] = *downloadPath
		anyFlagUsed = true
	}

	if *rateLimit != "" {
		flagsUsed["rate-limit"] = *rateLimit
		anyFlagUsed = true
	}

	if !*logToFile {
		flagsUsed["B"] = "true"
		anyFlagUsed = true
	} else {
		flagsUsed["B"] = "false"
	}

	return flagsUsed, anyFlagUsed, *web, url
}
