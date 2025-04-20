package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"wget/utils"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/time/rate"
)

// HandleDownloadWithFlags manages downloading from URL with various CLI flags.
func HandleDownloadWithFlags(url string, flags map[string]string) {
	var (
		err                     error
		logToFile               bool
		saveInDifferentLocation bool
		changeFileName          bool
		fileName                string
		filePath                string
		joinedPath              string
		rateLimit               = -1
		logger                  *log.Logger
		logWriter               io.Writer
	)

	for key, value := range flags {
		switch key {
		case "B":
			logToFile = true
		case "O":
			changeFileName = true
			fileName = value
		case "P":
			saveInDifferentLocation = true
			filePath, err = expandPath(value)
			if err != nil {
				log.Fatalf("Error expanding path: %v", err)
			}
		case "i":
			inputFile := flags["i"]
			SetFileName(inputFile)
			FileList(inputFile)
			return
		case "rate-limit":
			rateLimit, err = adjustRateLimit(value)
			if err != nil {
				log.Fatalf("Error adjusting rate limit: %v", err)
			}
		}
	}

	if rateLimit > 0 {
		limiter = rate.NewLimiter(rate.Limit(rateLimit), 64*1024) // 64 KB burst
	}

	// Log writer setup
	if logToFile {
		logFile, err := os.OpenFile("wget-log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatalf("Error opening log file: %v", err)
		}
		defer logFile.Close()
		logWriter = logFile
	} else {
		logWriter = os.Stdout
	}

	logger = log.New(logWriter, "", log.LstdFlags)
	logger.Printf("start at %v", time.Now().Format("2006-01-02 15:04:05"))

	// Download request
	resp, err := http.Get(url)
	if err != nil {
		logger.Fatalf("Error downloading: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Fatalf("Status: %v", resp.Status)
	}
	logger.Println("Request successful - status 200 OK")

	size := resp.ContentLength
	logger.Printf("Content size: %d bytes (~%.2f MB)", size, float64(size)/(1024*1024))

	if !changeFileName {
		fileName, err = utils.MakeAName(url)
		if err != nil {
			logger.Fatalf("Error creating filename: %v", err)
		}
	}

	// Build file path
	if saveInDifferentLocation {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			logger.Fatalf("Failed to create directory %s: %v", filePath, err)
		}
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			logger.Fatalf("Error getting absolute path: %v", err)
		}
		joinedPath = filepath.Join(absPath, fileName)
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			logger.Fatalf("Error getting current working directory: %v", err)
		}
		joinedPath = filepath.Join(cwd, fileName)
	}

	file, err := os.Create(joinedPath)
	if err != nil {
		logger.Fatalf("Error creating file: %v", err)
	}
	defer file.Close()
	logger.Printf("Saving file to: %s", joinedPath)

	// Setup writer with optional progress bar
	var writer io.Writer
	if logToFile {
		writer = io.MultiWriter(file)
		fmt.Println("Output will be written to ‘wget-log’.")
	} else {
		bar := progressbar.DefaultBytes(size, "Downloading")
		writer = io.MultiWriter(file, bar)
	}

	reader := resp.Body
	if limiter != nil {
		reader = &rateLimitedReader{ReadCloser: resp.Body, limiter: limiter}
	}

	if _, err := io.Copy(writer, reader); err != nil {
		logger.Fatalf("Error writing to file: %v", err)
	}

	logger.Printf("Downloaded [%s] finished at %s", url, time.Now().Format("2006-01-02 15:04:05"))
}

// expandPath replaces ~ with the user's home directory
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