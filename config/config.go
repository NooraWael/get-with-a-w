package config

import (
	"flag"
	"os"
	"wget/logger"
	"wget/utils"
)

//configuration settings and command line flags
//here we call this when there are flags

//bery bery complicated but fun

func HandleDownloadWithFlags(url string, flags map[string]string) {
	// Here you would handle different flags and call the appropriate downloading logic
	for key, value := range flags {
		switch key {
		case "B":
			logger.DownloadAndLog(url, value)
		}
	}
}

func ParseFlags() (map[string]string, bool, bool,string) {
	// Define all possible flags
	outputFileName := flag.String("O", "", "Specify the output file name (optional)")
	downloadPath := flag.String("P", "", "Specify the path to save the file")
	downladAndLog := flag.String("B", "", "Specifiy the filename to write the log into")
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

	url := ""
	// Check each flag and add to the map if it was set
	if *outputFileName != "" {
		flagsUsed["O"] = *outputFileName
		url = *outputFileName
		anyFlagUsed = true
	}
	if *downloadPath != "" {
		flagsUsed["P"] = *downloadPath
		url = *downloadPath
		anyFlagUsed = true
	}
	if *rateLimit != "" {
		flagsUsed["rate-limit"] = *rateLimit
		url = *rateLimit
		anyFlagUsed = true
	}
	if *downladAndLog != "" {
		flagsUsed["B"] = *downladAndLog
		url = *downladAndLog
		anyFlagUsed = true
	}
	


	return flagsUsed, anyFlagUsed, *web,url
}
