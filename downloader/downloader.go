package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

var(
	fileName string
	outputPath string
	mirrorMode bool
	multiFileMode bool
)
// DownloadFile downloads a file from the specified URL to the given output path.
// It logs the start and end time of the download, checks the HTTP response status,
// and writes the file to the specified path with real-time progress updates.
//
// @param url - The URL from which the file will be downloaded.
// @param outputPath - The file system path where the downloaded file will be saved.
// @return error - Returns an error if the request fails, the response is not 200 OK,
//
//	the file cannot be created, or the write operation fails.
//
// @example
// // Example of downloading a file:
// err := DownloadFile("https://example.com/file.zip", "./file.zip")
//
//	if err != nil {
//	    fmt.Println("Download failed:", err)
//	}
func DownloadFile(url string) (*os.File, error) {
	startTime := time.Now()
	fmt.Printf("start at %s\n", startTime.Format("2006-01-02 15:04:05"))

	// Sending request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("sending request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %s", resp.Status)
	}
	fmt.Println("sending request, awaiting response... status 200 OK")

	// Get content size
	size := resp.ContentLength
	sizeMB := float64(size) / (1024 * 1024)
	fmt.Printf("content size: %d [~%.2fMB]\n", size, sizeMB)

	// Create the file
	file, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()
	fmt.Printf("saving file to: %s\n", fileName)

	// Create progress bar
	bar := progressbar.DefaultBytes(
		size,
		"downloading",
	)

	// Create a multiwriter to write to both file and progress bar
	writer := io.MultiWriter(file, bar)

	// Copy the body to file and progress bar
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error writing to file: %v", err)
	}

	finishTime := time.Now()
	fmt.Printf("\nDownloaded [%s]\nfinished at %s\n", url, finishTime.Format("2006-01-02 15:04:05"))
	return file, nil
}

func SetFileName(name string) {
	fileName = name
}

func GetFileName() string {
	return fileName
}

func GetOutputPath() string {
	return outputPath
}

func SetOutputPath(path string) {
	outputPath = path
}

func SetMirrorMode(mode bool) {
	mirrorMode = mode
}

func SetMultiFileMode(mode bool) {
	multiFileMode = mode
}