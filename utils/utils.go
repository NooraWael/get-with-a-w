package utils

import (
	"net/url"
	"path"
	"strings"
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
