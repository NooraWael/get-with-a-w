package web

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"wget/downloader"

	"github.com/gin-gonic/gin"
)

func StartWebServer() {
	router := gin.Default()

	// Serve static files
	router.Static("/static", "./web/static")

	// Load templates
	router.LoadHTMLGlob("web/templates/*")

	// Define routes
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	router.GET("/documentation", func(c *gin.Context) {
		c.HTML(http.StatusOK, "documentation.html", nil)
	})

	router.POST("/download", func(c *gin.Context) {
		url := c.PostForm("url")
		downloadDirectory := getDownloadsPath()
		filename := filepath.Base(url) // Extract the filename from the URL

		// Complete path where the file will be saved
		outputPath := filepath.Join(downloadDirectory, filename)

		// Ensure the download directory exists
		if err := os.MkdirAll(downloadDirectory, 0755); err != nil {
			c.String(http.StatusInternalServerError, "Failed to create download directory: %s", err.Error())
			return
		}

		// Call your existing download function
		_, err := downloader.DownloadFile(url, false)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to download: %s", err.Error())
			return
		}
		c.String(http.StatusOK, "File downloaded successfully to %s", outputPath)
	})

	// Start the server
	router.Run(":8080")
}

func getDownloadsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return "./" // Fallback to current directory if the home directory can't be determined
	}
	// Directly join the Downloads directory to the home directory
	return filepath.Join(home, "Downloads")
}
