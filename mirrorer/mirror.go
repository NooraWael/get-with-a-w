package mirrorer

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"wget/downloader"

	"github.com/PuerkitoBio/goquery"
)

var excludeExtsList = []string{}
var excludeDirsList = []string{}
var convertLinks = false
var baseURL *url.URL

func DownloaderWrapper(urlStr string) *os.File {
	file, err := downloader.DownloadFile(urlStr)
	if err != nil {
		fmt.Printf("Error downloading file: %v\n", err)
		return nil
	}
	if file == nil {
		fmt.Printf("Error: file is nil\n")
		return nil
	} else {
		return file
	}
}

// SetExcludeExtsList sets the list of file extensions to exclude during mirroring.
func SetExcludeExtsList(list []string) {
	excludeExtsList = list
}

// SetExcludeDirsList sets the list of directories to exclude during mirroring.
func SetExcludeDirsList(list []string) {
	excludeDirsList = list
}

// SetConvertLinks sets whether to convert links during mirroring.
func SetConvertLinks(convert bool) {
	convertLinks = convert
}

// Mirror downloads the specified URL and mirrors its contents.
func Mirror(url *url.URL) {
	baseURL = url
	file := DownloaderWrapper(url.String())
	if file == nil {
		return
	}
	defer file.Close()

	file.Seek(0, 0)
	patchLinks(file)
}

func processUrl(urlStr string) string {
	// Parse the base URL if it's not already parsed
	if baseURL == nil {
		return urlStr
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}

	// If the input URL is absolute, return it directly
	if u.IsAbs() {
		return urlStr
	}

	if urlStr == "" {
		return baseURL.String()
	}

	// If the input URL is a relative path starting with "/", resolve it against the baseURL
	if strings.HasPrefix(urlStr, "/") {
		// Join base path and urlStr
		return fmt.Sprintf("%s://%s%s", baseURL.Scheme, baseURL.Host, path.Join(baseURL.Path, urlStr))
	}

	if baseURL.Path == "/" {
		return fmt.Sprintf("%s://%s/%s", baseURL.Scheme, baseURL.Host, urlStr)
	}

	// Otherwise, treat it as a relative path and resolve it against the baseURL's path
	return fmt.Sprintf("%s://%s%s/%s", baseURL.Scheme, baseURL.Host, baseURL.Path, urlStr)
}

func linkAllowed(link string) bool {
	for _, ext := range excludeExtsList {
		if strings.HasSuffix(link, "."+ext) {
			return false
		}
	}

	for _, dir := range excludeDirsList {
		if strings.HasPrefix(strings.TrimPrefix(link, "."), dir) {
			return false
		}
	}

	return true
}

func patchLinks(file *os.File) {
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return
	}

	wg := sync.WaitGroup{}

	// extract A tag href attributes
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		if linkAllowed(link) {
			wg.Add(1)
			go func() {
				defer wg.Done()
				linkFile := DownloaderWrapper(processUrl(link))
				if linkFile == nil {
					return
				}
				defer linkFile.Close()

				if convertLinks {
					relativePath, err := filepath.Rel(filepath.Dir(file.Name()), linkFile.Name())

					if err != nil {
						s.SetAttr("href", linkFile.Name())
						return
					}
					s.SetAttr("href", relativePath)
				}
			}()
		}
	})

	// extract IMG tag src attributes
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("src")

		if linkAllowed(link) {
			wg.Add(1)
			go func() {
				defer wg.Done()
				linkFile := DownloaderWrapper(processUrl(link))
				if linkFile == nil {
					return
				}
				defer linkFile.Close()

				if convertLinks {

					relativePath, err := filepath.Rel(filepath.Dir(file.Name()), linkFile.Name())

					if err != nil {
						s.SetAttr("src", linkFile.Name())
						return
					}
					s.SetAttr("src", relativePath)
				}
			}()
		}
	})

	// extract LINK tag href attributes
	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		if linkAllowed(link) {
			wg.Add(1)
			go func() {
				defer wg.Done()
				linkFile := DownloaderWrapper(processUrl(link))
				if linkFile == nil {
					return
				}
				defer linkFile.Close()

				if convertLinks {

					relativePath, err := filepath.Rel(filepath.Dir(file.Name()), linkFile.Name())

					if err != nil {
						s.SetAttr("href", linkFile.Name())
						return
					}
					s.SetAttr("href", relativePath)
				}
			}()
		}
	})

	// extract SCRIPT tag src attributes
	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("src")

		if linkAllowed(link) {
			wg.Add(1)
			go func() {
				defer wg.Done()
				linkFile := DownloaderWrapper(processUrl(link))
				if linkFile == nil {
					return
				}
				defer linkFile.Close()

				if convertLinks {

					relativePath, err := filepath.Rel(filepath.Dir(file.Name()), linkFile.Name())

					if err != nil {
						s.SetAttr("src", linkFile.Name())
						return
					}
					s.SetAttr("src", relativePath)
				}
			}()
		}
	})

	// internal CSS links
	doc.Find("style").Each(func(i int, s *goquery.Selection) {
		urlMatcher := regexp.MustCompile(`url\(['"]?(.*?)['"]?\)`)

		css := s.Text()

		css = urlMatcher.ReplaceAllStringFunc(css, func(match string) string {
			url := urlMatcher.FindStringSubmatch(match)[1]

			if !linkAllowed(url) {
				return match
			}

			linkFile := DownloaderWrapper(processUrl(url))
			if linkFile == nil {
				return match
			}
			defer linkFile.Close()

			if convertLinks {

				relativePath, err := filepath.Rel(filepath.Dir(file.Name()), linkFile.Name())

				if err != nil {
					return match
				}

				return strings.Replace(match, url, relativePath, 1)
			}
			return url
		})
		s.SetHtml(css)
	})

	wg.Wait()


	// save the modified HTML
	docHtml, err := doc.Html()
	if err != nil {
		return
	}

	file.Seek(0, 0)
	file.Truncate(0)
	file.WriteString(docHtml)

}
