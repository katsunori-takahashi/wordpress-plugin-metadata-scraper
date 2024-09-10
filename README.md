# WordPress Plugin Version Crawler

This tool is designed to scrape metadata from WordPress plugin pages and export the information to a CSV file.

## Features

- Reads a list of WordPress plugin URLs from a CSV file
- Scrapes the following metadata for each plugin:
  - Plugin Name
  - Version
  - Last Updated Date
  - Active Installations
  - Required WordPress Version
  - Tested Up To Version
  - Required PHP Version
  - Supported Languages
  - Tags
- Implements retry logic for handling rate limiting (HTTP 429 errors)
- Exports collected data to a CSV file
- Logs all operations for easy debugging and monitoring

## How it works

1. The program reads plugin URLs from a CSV file named `plugin_urls.csv`.
2. It then visits each URL and scrapes the relevant metadata.
3. If a rate limit error occurs, the program will wait and retry the request.
4. All scraped data is collected and exported to a file named `plugin_meta_results.csv`.
5. The entire process is logged to `scraper.log` for monitoring and debugging purposes.

## Usage

1. Ensure you have a `plugin_urls.csv` file with a list of WordPress plugin URLs.
2. Run the program: `go run main.go`
3. Check the `plugin_meta_results.csv` for the scraped data and `scraper.log` for the operation log.

Note: This tool is designed for educational and research purposes. Please respect WordPress.org's terms of service and rate limiting policies when using this tool.
