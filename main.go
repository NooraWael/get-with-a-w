package main

import (
	"flag"
	"fmt"
	"wget/config"
	"wget/downloader"
	"wget/utils"
)

func main() {
	//get the flags entered in
	flags, flagProvided := config.ParseFlags()

	url := flag.Arg(0)
	if flagProvided {
		config.HandleDownloadWithFlags(url, flags)
	} else {
		fmt.Println("No flags provided.")

		//get a name for the download and call the download function
		output,err := utils.MakeAName(url) 
		if err != nil {
			fmt.Println("Error making a name for the download:", err)
			return
		}

		//make sure there is :// or something like that
		url = utils.EnsureScheme(url)
		err = downloader.DownloadFile(url, output)
		if err != nil {
			fmt.Println("Error downloading the file:", err)
		}
	}

}