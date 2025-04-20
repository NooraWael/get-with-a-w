package mirrorer

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"flag"
)
// code when mirror flag is set
func ParseMirrorFlag(flags map[string]string) {
	if flags["reject"] != "" {
		flags["R"] = flags["reject"]
	}

	if flags["exclude"] != "" {
		flags["X"] = flags["exclude"]
	}


	if flags["R"] != "" {
		rejectList := strings.Split(flags["R"], ",")
		SetExcludeExtsList(rejectList)
	}

	if flags["X"] != "" {
		excludeDirs := strings.Split(flags["X"], ",")
		SetExcludeDirsList(excludeDirs)
	}

	if flags["convertLinks"] != "" {
		SetConvertLinks(true)
	}

	if flag.NArg() == 0 {
		fmt.Println("Missing URL")
		os.Exit(1)
	}

	url, _ := url.Parse(flag.Arg(0))
	println("Mirroring URL:", url.String())
	Mirror(url)
}
