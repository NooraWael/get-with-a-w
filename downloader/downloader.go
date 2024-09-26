package downloader

import (
	"net/http"
	"os"
)

//manages the downloading files logic

// DownloadFile downloads a file from the specified URL to the given output path.
func DownloadFile(url, outputPath string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Create the file at the specified path
    file, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer file.Close()

    // Write the body to file
    _, err = file.ReadFrom(resp.Body)
    if err != nil {
        return err
    }

    return nil
}