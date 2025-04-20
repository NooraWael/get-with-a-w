package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
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
				mirror(flags)
				return
			} else {
				config.HandleDownloadWithFlags(url2, flags)
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

// code when mirror flag is set
func mirror(flags map[string]string) {
	if flags["reject"] != "" {
		flags["R"] = flags["reject"]
	}

	if flags["exclude"] != "" {
		flags["X"] = flags["exclude"]
	}


	if flags["R"] != "" {
		rejectList := strings.Split(flags["R"], ",")
		mirrorer.SetExcludeExtsList(rejectList)
	}

	if flags["X"] != "" {
		excludeDirs := strings.Split(flags["X"], ",")
		mirrorer.SetExcludeDirsList(excludeDirs)
	}

	if flags["convertLinks"] != "" {
		mirrorer.SetConvertLinks(true)
	}

	if flag.NArg() == 0 {
		fmt.Println("Missing URL")
		os.Exit(1)
	}

	url, _ := url.Parse(flag.Arg(0))
	println("Mirroring URL:", url.String())
	mirrorer.Mirror(url)
}
