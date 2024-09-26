package config

import (
	"flag"
	"fmt"
)

//configuration settings and command line flags
//here we call this when there are flags

//bery bery complicated but fun

func HandleDownloadWithFlags(url string, flags map[string]string) {
    // Here you would handle different flags and call the appropriate downloading logic
    fmt.Println("Downloading with flags:", flags)
   
}

func ParseFlags() (map[string]string, bool) {
    outputFileName := flag.String("O", "", "Specify the output file name (optional)")
    flag.Parse()

    flagsUsed := make(map[string]string)
    anyFlagUsed := false

    if *outputFileName != "" {
        flagsUsed["O"] = *outputFileName
        anyFlagUsed = true
    }

    return flagsUsed, anyFlagUsed
}