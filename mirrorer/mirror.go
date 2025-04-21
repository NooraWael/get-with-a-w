// Package mirrorer handles website mirroring by downloading resources and updating links
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

var (
	excludeExtsList = []string{}
	excludeDirsList = []string{}
	convertLinks    = false
	baseURL         *url.URL
)

// DownloaderWrapper wraps the file download logic for reuse.
func DownloaderWrapper(urlStr string) *os.File {
	file, err := downloader.DownloadFile(urlStr, true)
	if err != nil {
		fmt.Printf("Error downloading file: %v\n", err)
		return nil
	}
	return file
}

// SetExcludeExtsList sets excluded file extensions.
func SetExcludeExtsList(list []string) {
	excludeExtsList = list
}

// SetExcludeDirsList sets excluded directory prefixes.
func SetExcludeDirsList(list []string) {
	excludeDirsList = list
}

// SetConvertLinks determines whether to convert absolute URLs to relative.
func SetConvertLinks(convert bool) {
	convertLinks = convert
}

// Mirror downloads and patches a given HTML page and its assets.
func Mirror(url *url.URL) {
	baseURL = url
	file := DownloaderWrapper(url.String())
	if file == nil {
		return
	}
	defer file.Close()
	file.Seek(0, 0)
	patchLinks(file)
	fmt.Println()
}

// processUrl resolves relative URLs using baseURL.
func processUrl(urlStr string) string {
	if baseURL == nil || urlStr == "" {
		return urlStr
	}

	u, err := url.Parse(urlStr)
	if err != nil || u.IsAbs() {
		return urlStr
	}

	if strings.HasPrefix(urlStr, "/") || strings.HasPrefix(urlStr, "./") {
		return fmt.Sprintf("%s://%s%s", baseURL.Scheme, baseURL.Host, path.Join(baseURL.Path, urlStr))
	}

	if baseURL.Path == "/" {
		return fmt.Sprintf("%s://%s/%s", baseURL.Scheme, baseURL.Host, urlStr)
	}

	return fmt.Sprintf("%s://%s%s/%s", baseURL.Scheme, baseURL.Host, baseURL.Path, urlStr)
}

// linkAllowed checks if a link should be excluded.
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

// patchLinks updates resource links in the HTML and downloads them.
func patchLinks(file *os.File) {
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return
	}

	wg := sync.WaitGroup{}

	downloadAndPatch := func(sel *goquery.Selection, attr string) {
		link, exists := sel.Attr(attr)
		if !exists || !linkAllowed(link) {
			return
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			linkFile := DownloaderWrapper(processUrl(link))
			if linkFile == nil {
				return
			}
			defer linkFile.Close()
			if convertLinks {
				relPath, err := filepath.Rel(filepath.Dir(file.Name()), linkFile.Name())
				if err != nil {
					sel.SetAttr(attr, linkFile.Name())
				} else {
					sel.SetAttr(attr, relPath)
				}
			}
		}()
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) { downloadAndPatch(s, "href") })
	doc.Find("img").Each(func(i int, s *goquery.Selection) { downloadAndPatch(s, "src") })
	doc.Find("link").Each(func(i int, s *goquery.Selection) { downloadAndPatch(s, "href") })
	doc.Find("script").Each(func(i int, s *goquery.Selection) { downloadAndPatch(s, "src") })

	// Handle inline CSS background images
	doc.Find("style").Each(func(i int, s *goquery.Selection) {
		matcher := regexp.MustCompile(`url\(['"]?(.*?)['"]?\)`)
		css := s.Text()
		css = matcher.ReplaceAllStringFunc(css, func(match string) string {
			url := matcher.FindStringSubmatch(match)[1]
			if !linkAllowed("url") {
				return match
			}
			linkFile := DownloaderWrapper(processUrl(url))
			if linkFile == nil {
				return match
			}
			defer linkFile.Close()
			if convertLinks {
				relPath, err := filepath.Rel(filepath.Dir(file.Name()), linkFile.Name())
				if err == nil {
					return strings.Replace(match, url, relPath, 1)
				}
			}
			return match
		})
		s.SetHtml(css)
	})

	wg.Wait()

	// Save updated HTML back to file
	html, err := doc.Html()
	if err != nil {
		return
	}
	file.Seek(0, 0)
	file.Truncate(0)
	file.WriteString(html)
}
