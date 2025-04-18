package main

import (
	"flag"
	"fmt"
	"os"
	"net/url"
	"wget/config"
	"wget/downloader"
	"wget/utils"
	"wget/web"
	"wget/mirrorer"
)

func main() {
	//get the flags entered in
	flags, flagProvided, startweb,url2 := config.ParseFlags()
	if startweb {
		web.StartWebServer()
	} else {
		url := flag.Arg(0)
		print(flags["i"], flagProvided)
		if flagProvided {
			if flags["i"] != "" {
				print(flags)
				inputFile := flags["i"]
				downloader.SetFileName(inputFile)
				downloader.SetMultiFileMode(true)
				downloader.FileList(inputFile)	
				return			
			}
			config.HandleDownloadWithFlags(url2, flags)
		} else {
			//get a name for the download and call the download function
			output, err := utils.MakeAName(url)
			if err != nil {
				fmt.Println("Error making a name for the download:", err)
				return
			}
			downloader.SetFileName(output)
			downloader.SetOutputPath(output)
			//make sure there is :// or something like that
			url = utils.EnsureScheme(url)
			_, err = downloader.DownloadFile(url)
			if err != nil {
				fmt.Println("Error downloading the file:", err)
			}
		}
	}
}

// code when mirror flag is set
func mirror(flags map[string]any) {
	if flags["convert-links"] != nil {
		convertLinks := flags["convert-links"].(*bool)
		mirrorer.SetConvertLinks(*convertLinks)
	}

	if flags["outputFileName"] != "" {
		fmt.Println("Cannot specify both -O and -mirror")
		os.Exit(1)
	}

	if flags["inputFile"]!= "" {
		fmt.Println("Cannot specify both -i and -mirror")
		os.Exit(1)
	}

	if flags["X"] != "" && flags["exclude"] != "" {
		fmt.Println("Cannot specify both -X and -exclude")
		os.Exit(1)
	}

	if flags["R"] != "" && flags["reject"] != "" {
		fmt.Println("Cannot specify both -R and -reject")
		os.Exit(1)
	}

	if flags["reject"] != "" {
		flags["R"] = flags["reject"]
	}

	if flags["exclude"] != "" {
		flags["X"] = flags["exclude"]
	}


	if flags["R"] != "" {
		rejectList := flags["R"].([]string)
		mirrorer.SetExcludeDirsList(rejectList)
	}

	if flags["X"] != "" {
		excludeDirs := flags["X"].([]string)
		mirrorer.SetExcludeDirsList(excludeDirs)
	}

	if flag.NArg() == 0 {
		fmt.Println("Missing URL")
		os.Exit(1)
	}

	downloader.SetMirrorMode(true)

	url, _ := url.Parse(flag.Arg(0))

	mirrorer.Mirror(url)
}
