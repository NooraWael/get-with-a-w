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
func DownloadFile(fileURL string, mirrorMode bool) (*os.File, error) {
	startTime := time.Now()
	fmt.Printf("start at %s\n", startTime.Format("2006-01-02 15:04:05"))
	// Sending request
	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, fmt.Errorf("sending request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}
	fmt.Println("sending request, awaiting response... status 200 OK")

	// Get content size
	size := resp.ContentLength
	sizeMB := float64(size) / (1024 * 1024)
	fmt.Printf("content size: %d [~%.2fMB]\n", size, sizeMB)
	if mirrorMode {
		parsedURL, err := url.Parse(fileURL)
		if err != nil {
			fmt.Printf("Error parsing URL: %v", err)
			return nil, err
		}
		if parsedURL.Path == "" {
			downloadPath = parsedURL.Host + string(os.PathSeparator)
		} else {
			downloadPath, _ = path.Split(parsedURL.Host + parsedURL.Path)
		}
		// Ensure the directory exists
		if err := os.MkdirAll(downloadPath, os.ModePerm); err != nil {
			fmt.Printf("Failed to create directory: %s, error: %v\n", downloadPath, err)
			return nil, nil
		}
	}
	fileName, err = utils.MakeAName(fileURL)
	if err != nil {
		fmt.Println("Error making a name for the download:", err)
		return nil, err
	}
	// Create the file
	if mirrorMode {
		fileName = filepath.Join(downloadPath, fileName)
	}
	file, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %v", err)
	}
	filePath := filepath.Dir(file.Name())
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatalf("Error getting path: %v", err)
	}

	// join the path of the folder to save the file into with the file name
	joinedPath := filepath.Join(absFilePath)
	fmt.Printf("saving file to: %s\n", joinedPath)
	fmt.Println("file name:", fileName)

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
	fmt.Printf("\nDownloaded [%s]\nfinished at %s\n", fileURL, finishTime.Format("2006-01-02 15:04:05"))
	return file, nil
}

func SetFileName(name string) {
	fileName = name
}
