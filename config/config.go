package config

import (
	"flag"
	"fmt"
	"wget/utils"
	"os"
)

//configuration settings and command line flags
//here we call this when there are flags

//bery bery complicated but fun

func HandleDownloadWithFlags(url string, flags map[string]string) {
    // Here you would handle different flags and call the appropriate downloading logic
    fmt.Println("Downloading with flags:", flags)
   
}
func ParseFlags() (map[string]string, bool) {
    // Define all possible flags
    outputFileName := flag.String("O", "", "Specify the output file name (optional)")
    downloadPath := flag.String("P", "", "Specify the path to save the file")
    rateLimit := flag.String("rate-limit", "", "Specify the maximum download rate (e.g., '500k', '2M')")
    help := flag.Bool("help", false, "Display help information")

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

    // Check each flag and add to the map if it was set
    if *outputFileName != "" {
        flagsUsed["O"] = *outputFileName
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

    return flagsUsed, anyFlagUsed
}