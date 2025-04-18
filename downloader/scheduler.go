package downloader

import (
	"fmt"
	"os"
	"sync"
	"bufio"
	"strings"
	// "wget/logger"
)
// Manages multiple and background downloads
// i flag mostly
// Handles the case when -i flag is set
func FileList(inputFile string) {
	SetMultiFileMode(true)
	// convert interface to string
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// WaitGroup to wait for all the goroutines to finish
	var wg sync.WaitGroup

	// Create a new buffered reader
	reader := bufio.NewReader(file)

	for {

		line, _, err := reader.ReadLine()
		if err != nil {
			if err.Error() != "EOF" {
				msg := fmt.Sprintf("Failed to read file: %v", err)
				print(msg)
				// logger.Log(msg)
			}
			break
		}

		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			print("Testing link: ", link)
			DownloadFile(strings.TrimSuffix(link, "\n"))
		}(string(line))
	}

	wg.Wait()
}
