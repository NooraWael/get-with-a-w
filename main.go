package main

import (
	"flag"
	"fmt"
	"wget/config"
	"wget/downloader"
	"wget/mirrorer"
	"wget/utils"
	"wget/web"
)

func main() {
	//get the flags entered in
	flags, flagProvided, startweb, url2, err := config.ParseFlags()
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		return
	}
	if startweb {
		web.StartWebServer()
	} else {
		url := flag.Arg(0)
		if flagProvided {
			if flags["mirror"] != "" {
				mirrorer.ParseMirrorFlag(flags)
				return
			} else {
				downloader.HandleDownloadWithFlags(url2, flags)
			}
		} else {
			//get a name for the download and call the download function
			output, err := utils.MakeAName(url)
			if err != nil {
				fmt.Println("Error making a name for the download:", err)
				return
			}
			downloader.SetFileName(output)

			//make sure there is :// or something like that
			url = utils.EnsureScheme(url)
			_, err = downloader.DownloadFile(url, false)
			if err != nil {
				fmt.Println("Error downloading the file:", err)
			}
		}
	}
}