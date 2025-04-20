// Package config handles CLI flags and the main download logic
package config

import (
	"flag"
	"fmt"
	"os"
	"wget/utils"
)

// ParseFlags parses command line arguments into a flag map
func ParseFlags() (map[string]string, bool, bool, string, error) {
	flagSet := map[string]*string{
		"O":        flag.String("O", "", "Specify the output file name (optional)"),
		"P":        flag.String("P", "", "Specify the path to save the file"),
		"i":        flag.String("i", "", "Specify the input file containing URLs"),
		"rate-limit": flag.String("rate-limit", "", "Specify the max download rate e.g. '500k', '2M'"),
		"R":        flag.String("R", "", "Comma separated list of file extensions to reject"),
		"reject":   flag.String("reject", "", "Alias for -R"),
		"X":        flag.String("X", "", "Comma separated list of directories to exclude"),
		"exclude":  flag.String("exclude", "", "Alias for -X"),
	}
	flagB := flag.Bool("B", false, "Log output to wget-log")
	flagMirror := flag.Bool("mirror", false, "Mirror the entire website")
	flagHelp := flag.Bool("help", false, "Display help information")
	flagWeb := flag.Bool("web", false, "Start the web server interface")
	flagConvert := flag.Bool("convert-links", false, "Convert links to local")

	flag.Parse()

	if *flagHelp {
		utils.DisplayHelp()
		os.Exit(0)
	}

	flagsUsed := make(map[string]string)
	urlArg := flag.Arg(0)
	anyUsed := false

	for key, val := range flagSet {
		if *val != "" {
			flagsUsed[key] = *val
			anyUsed = true
		}
	}

	if *flagB {
		flagsUsed["B"] = "wget-log"
		anyUsed = true
	}
	if *flagMirror {
		flagsUsed["mirror"] = "mirror"
		anyUsed = true
	}
	if *flagConvert {
		flagsUsed["convertLinks"] = "true"
		anyUsed = true
	}

	// validation for mutually exclusive flags
	conflicts := [][2]string{
		{"i", "O"}, {"i", "P"}, {"i", "B"}, {"i", "rate-limit"},
		{"R", "reject"}, {"X", "exclude"}, {"mirror", "O"},
		{"mirror", "i"}, {"mirror", "P"}, {"mirror", "B"},
		{"mirror", "rate-limit"},
	}
	for _, pair := range conflicts {
		if flagsUsed[pair[0]] != "" && flagsUsed[pair[1]] != "" {
			return nil, false, false, "", fmt.Errorf("cannot specify both -%s and -%s", pair[0], pair[1])
		}
	}

	if (flagsUsed["R"] != "" || flagsUsed["reject"] != "") &&
		(flagsUsed["X"] != "" || flagsUsed["exclude"] != "") &&
		flagsUsed["mirror"] == "" {
		return nil, false, false, "", fmt.Errorf("cannot use -reject or -exclude without -mirror")
	}

	return flagsUsed, anyUsed, *flagWeb, urlArg, nil
}
