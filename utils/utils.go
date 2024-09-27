package utils

import (
	"net/url"
	"path"
	"strings"
	"fmt"
)

//any utility functions that will optimize the program

func MakeAName(urlStr string) (string,error) {

	//parse the url 
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return "default",err
	}

	//Extract base
	fileName := path.Base(parsedUrl.Path)

	if fileName == "" || !strings.Contains(fileName, ".") {
        return "default_filename", nil
    }

    return fileName,nil
}

func EnsureScheme(urlStr string) string {
    if !strings.Contains(urlStr, "://") {
        urlStr = "https://" + urlStr // Default to using HTTPS if no scheme is provided
    }
    return urlStr
}


func DisplayHelp() {
    fmt.Println(`Usage: go run . [options] <URL>
Options:
  -B                  Run download in background and output to 'wget-log'.
  -O <filename>       Download as a different filename.
  -P <path>           Path where the file will be saved.
  --rate-limit <rate> Limit the download rate (e.g., 500k, 2M).
  -i <file>           Download multiple files listed in a file.
  --mirror            Download an entire website for offline viewing.
  -R <types>          Reject files of specified types (e.g., jpg, gif), used with --mirror.
  -X <paths>          Exclude certain paths from being downloaded, used with --mirror.
  --convert-links     Convert links for offline viewing, used with --mirror.

Examples:
  go run . https://example.com/file.zip
  go run . -O myfile.zip https://example.com/file.zip
  go run . --rate-limit=1M https://example.com/bigfile.zip
  go run . --mirror --convert-links https://example.com

Use 'man wget' for more information on wget features.`)
}