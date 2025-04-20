// Package downloader handles file downloading logic
package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"wget/utils"

	"github.com/schollz/progressbar/v3"
)

var (
	fileName      string
	downloadPath  string
	mirrorMode    bool
	multiFileMode bool
)

// DownloadFile downloads a file from the specified URL and saves it locally.
// If mirrorMode is enabled, it preserves the directory structure from the URL.
// Displays a progress bar during the download.
func DownloadFile(fileURL string, mirrorMode bool) (*os.File, error) {
	startTime := time.Now()
	fmt.Printf("Start at %s\n", startTime.Format("2006-01-02 15:04:05"))

	// Perform HTTP GET request
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("sending request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil // No file downloaded for non-OK status
	}
	fmt.Println("Sending request, awaiting response... status 200 OK")

	// Log content size
	size := resp.ContentLength
	fmt.Printf("Content size: %d [~%.2fMB]\n", size, float64(size)/(1024*1024))

	// Generate target download path based on mirror mode
	if mirrorMode {
		parsedURL, err := url.Parse(fileURL)
		if err != nil {
			return nil, fmt.Errorf("error parsing URL: %v", err)
		}

		if parsedURL.Path == "" {
			downloadPath = parsedURL.Host + string(os.PathSeparator)
		} else {
			downloadPath, _ = path.Split(parsedURL.Host + parsedURL.Path)
		}

		// Ensure download directory exists
		if err := os.MkdirAll(downloadPath, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create directory: %s, error: %v", downloadPath, err)
		}
	}

	// Generate a file name for the downloaded content
	fileName, err = utils.MakeAName(fileURL)
	if err != nil {
		return nil, fmt.Errorf("error generating file name: %v", err)
	}

	if mirrorMode {
		fileName = filepath.Join(downloadPath, fileName)
	}

	// Create destination file
	file, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %v", err)
	}

	// Print download destination
	absFilePath, err := filepath.Abs(filepath.Dir(file.Name()))
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}
	fmt.Printf("Saving file to: %s\n", filepath.Join(absFilePath))
	fmt.Println("File name:", fileName)

	// Setup progress bar and multi-writer
	bar := progressbar.DefaultBytes(size, "Downloading")
	writer := io.MultiWriter(file, bar)

	// Perform the file download
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error writing to file: %v", err)
	}

	finishTime := time.Now()
	fmt.Printf("\nDownloaded [%s]\nFinished at %s\n", fileURL, finishTime.Format("2006-01-02 15:04:05"))
	return file, nil
}

// SetFileName allows setting the output file name manually.
func SetFileName(name string) {
	fileName = name
}
