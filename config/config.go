package config

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"wget/downloader"
	// "wget/mirrorer"
	"wget/utils"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/time/rate"
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
	logToFile := false               // check the existance of the -B flag
	saveInDifferentLocation := false // check the existance of the -P flag

	// ----------- -O flag -------------
	changeFileName := false // if the -O flag is passed change the file name that will be created
	var fileName string     // store the value of the file name

	// ----------- -P flag -------------
	var filePath string   // the file path that the file will be stored in using the value in the flag
	var joinedPath string // the entire path that will include the location and the filename

	var logger *log.Logger
	var logWriter io.Writer // variable that will handle where the log will be printed to
	var rateLimit int = -1
	var rateLimitErr error

	for key, value := range flags {
		switch key {
		case "B":
			logToFile = true
			fmt.Println("Output will be written to wget-log.")
		case "O":
			changeFileName = true
			fileName = value
		case "P":
			saveInDifferentLocation = true
			filePath, err = expandPath(value)
			if err != nil {
				log.Fatalf("Error expanding path: %v", err)
				os.Exit(1)
			}
		case "i":
			inputFile := flags["i"]
			downloader.SetFileName(inputFile)
			downloader.SetMultiFileMode(true)
			downloader.FileList(inputFile)
			return
		case "rate-limit":
			rateLimit, rateLimitErr = adjust_rate_limit(value)
			if rateLimitErr != nil {
				log.Fatalf("Error adjusting rate limit: %v", rateLimitErr)
				os.Exit(1)
			}
		case "mirror":
			// mirrorer
		}
	}
	// rateLimit is in bytes per second
	var limiter *rate.Limiter
	if rateLimit > 0 {
		limiter = rate.NewLimiter(rate.Limit(rateLimit), 64*1024) // burst of 64 KB
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

	if saveInDifferentLocation {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			fmt.Println("Failed to create directory: %s, error: %v\n", filePath, err)
			return
		}
		// get the absoulute path
		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			logger.Fatalf("Error getting path: %v", err)
		}

		// join the path of the folder to save the file into with the file name
		joinedPath = filepath.Join(absFilePath, fileName)
	} else {
		// get the absoulute path
		absFilePath, err := os.Getwd()
		if err != nil {
			logger.Fatalf("Error getting path: %v", err)
		}

		// join the path of the folder to save the file into with the file name
		joinedPath = filepath.Join(absFilePath, fileName)
	}

	file, err := os.Create(joinedPath) // Always save downloaded file as 'test'
	if err != nil {
		logger.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()
	logger.Printf("saving file to: ./%s", joinedPath)

	// Create a progress bar that writes to discard since we don't need to display it
	if !logToFile {
	bar := progressbar.DefaultBytes(size, "downloading")
	
	multiWriter := io.MultiWriter(file, bar)
	reader := resp.Body
	
	if limiter != nil {
		reader = &rateLimitedReader{
			ReadCloser: resp.Body,
			limiter:    limiter,
		}
	}
	_, err = io.Copy(multiWriter, reader)
	println(joinedPath)
}
	
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
	logToFile := flag.Bool("B", false, "Specify the filename to write the log into")
	inputFile := flag.String("i", "", "Specify the input file containing URLs")
	rateLimit := flag.String("rate-limit", "", "Specify the maximum download rate (e.g., '500k', '2M')")
	mirror := flag.Bool("mirror", false, "Mirror the entire website")
	help := flag.Bool("help", false, "Display help information")
	web := flag.Bool("web", false, "Start the web server interface")
	rejectList := flag.String("R", "", "Comma separated list of file extensions to reject")
	rejectListAlias := flag.String("reject", "", "Comma separated list of file extensions to reject (alias for -R)")

	excludeDirs := flag.String("X", "", "Comma separated list of directories to exclude")
	excludeDirsAlias := flag.String("exclude", "", "Comma separated list of directories to exclude (alias for -X)")
	convertLinks := flag.Bool("convert-links", false, "Convert links to local")
	// Parse the command line arguments
	flag.Parse()

	// Check if help was requested
	if *help {
		utils.DisplayHelp()
		os.Exit(0)
	}

	// Map to store flags that were actually set
	flagsUsed := make(map[string]string) // using generics to store multiple type at the same time
	anyFlagUsed := false

	// captures the non flag arguements
	url := flag.Arg(0)

	// Check each flag and add to the map if it was set
	if *outputFileName != "" {
		flagsUsed["O"] = *outputFileName
		anyFlagUsed = true
	}

	if *mirror {
		flagsUsed["mirror"] = "mirror"
		anyFlagUsed = true
	}

	if *rejectList != "" {
		flagsUsed["R"] = *rejectList
		anyFlagUsed = true
	}
	if *rejectListAlias != "" {
		flagsUsed["reject"] = *rejectListAlias
		anyFlagUsed = true
	}
	if *excludeDirs != "" {
		flagsUsed["X"] = *excludeDirs
		anyFlagUsed = true
	}

	if *excludeDirsAlias != "" {
		flagsUsed["exclude"] = *excludeDirsAlias
		anyFlagUsed = true
	}

	if *convertLinks != false {
		flagsUsed["convert-links"] = "convertLinks"
		anyFlagUsed = true
	}

	if *inputFile != "" {
		print(*inputFile)
		flagsUsed["i"] = *inputFile
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

	if *logToFile {
		flagsUsed["B"] = "wget-log"
		anyFlagUsed = true
	}

	return flagsUsed, anyFlagUsed, *web, url
}

func adjust_rate_limit(rateLimit string) (int, error) {
	// Parse the rate limit string and convert it to bytes per second
	var multiplier int
	testRateLimit, err := strconv.Atoi(rateLimit[:len(rateLimit)])
	if err == nil {
		return testRateLimit, nil
	}
	switch rateLimit[len(rateLimit)-1] {
	case 'k':
		multiplier = 1000
	case 'M':
		multiplier = 1000 * 1000
	case 'G':
		multiplier = 1000 * 1000 * 1000
	default:
		return 0, fmt.Errorf("invalid rate limit format")
	}

	rate, err := strconv.Atoi(rateLimit[:len(rateLimit)-1])
	if err != nil {
		return 0, err
	}

	return rate * multiplier, nil
}

type rateLimitedReader struct {
	io.ReadCloser
	limiter *rate.Limiter
}

func (r *rateLimitedReader) Read(p []byte) (int, error) {
	n, err := r.ReadCloser.Read(p)
	if err != nil {
		return n, err
	}
	// Wait for permission to read n bytes
	if err := r.limiter.WaitN(context.Background(), n); err != nil {
		return n, err
	}
	return n, nil
}

func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[2:]), nil
	}
	return path, nil
}
