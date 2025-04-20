# Get-With-A-W

Get-With-A-W is a versatile file downloader and website mirroring tool inspired by `wget`. It provides both a command-line interface (CLI) and a web-based interface for downloading files, mirroring websites, and managing downloads with advanced features like rate limiting, background logging, and more.

---

## Features

- **File Downloading**: Download files from URLs with support for custom filenames and save paths.
- **Website Mirroring**: Mirror entire websites for offline use, with options to exclude specific file types or directories.
- **Rate Limiting**: Control download speeds to avoid overloading networks.
- **Background Logging**: Log download activities to a file (`wget-log`) for later review.
- **Web Interface**: A user-friendly web interface for initiating downloads.
- **Collapsible Documentation**: Interactive documentation for easy navigation.
- **Progress Bar**: Visual feedback for download progress in the CLI.
- **Multi-File Downloads**: Download multiple files listed in a text file.
- **Link Conversion**: Convert links for offline viewing when mirroring websites.

---

## Project Structure

The project is divided into several modules for better organization:

- **`main.go`**: The entry point of the application. Handles CLI flags and determines whether to start the web server or perform a download.
- **`config`**: Handles parsing and validation of CLI flags.
- **`downloader`**: Contains the core logic for downloading files, handling flags, and managing rate limits.
- **`mirrorer`**: Handles website mirroring, including downloading resources and converting links.
- **`utils`**: Provides utility functions like URL validation, filename generation, and help display.
- **`logger`**: Manages logging to the console or a file.
- **`web`**: Implements the web server interface using the Gin framework.
- **`templates`**: HTML templates for the web interface.
- **`static`**: Static assets like CSS, JavaScript, and images for the web interface.

---

## Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-repo/get-with-a-w.git
   cd get-with-a-w
   ```

2. **Install Dependencies**:
   Ensure you have Go installed (version 1.23 or later). Run:
   ```bash
   go mod tidy
   ```

3. **Run the Application**:
   ```bash
   go run main.go
   ```

---

## Usage

### CLI Mode

Run the application with the following syntax:
```bash
go run main.go [options] <URL>
```

#### Options:
- `-B`: Log output to `wget-log`.
- `-O <filename>`: Save the file with a custom name.
- `-P <path>`: Specify the directory to save the file.
- `--rate-limit <rate>`: Limit download speed (e.g., `500k`, `2M`).
- `-i <file>`: Download multiple files listed in a text file.
- `--mirror`: Mirror an entire website.
- `-R <types>`: Reject files of specified types (e.g., `jpg`, `gif`).
- `-X <paths>`: Exclude specific paths from mirroring.
- `--convert-links`: Convert links for offline viewing.

#### Examples:
1. Download a single file:
   ```bash
   go run main.go https://example.com/file.zip
   ```

2. Download with a custom filename:
   ```bash
   go run main.go -O custom_name.zip https://example.com/file.zip
   ```

3. Mirror a website:
   ```bash
   go run main.go --mirror --convert-links https://example.com
   ```

4. Download multiple files from a list:
   ```bash
   go run main.go -i downloads.txt
   ```

5. Limit download speed:
   ```bash
   go run main.go --rate-limit=500k https://example.com/largefile.zip
   ```

### Web Interface

Start the web server:
```bash
go run main.go -web
```

Access the web interface at [http://localhost:8080](http://localhost:8080).

#### Features:
- **Home Page**: Enter a URL to download files.
- **Documentation Page**: Interactive documentation with collapsible sections for easy navigation.

---

## File Structure

```
get-with-a-w/
├── config/               # CLI flag parsing and validation
├── downloader/           # File downloading logic
├── logger/               # Logging functionality
├── mirrorer/             # Website mirroring logic
├── utils/                # Utility functions
├── web/                  # Web server implementation
│   ├── templates/        # HTML templates
│   └── static/           # Static assets (CSS, JS, images)
├── main.go               # Entry point
├── go.mod                # Go module dependencies
└── README.md             # Project documentation
```

## Authors

- **Noora Wael**: Main downloading function.
- **Ahmed AlAli**: Flags implementation and mirroring.
- **Hussain Almakana**: Flags implementation.
- **Hashem Alzaki**: Moral support.
