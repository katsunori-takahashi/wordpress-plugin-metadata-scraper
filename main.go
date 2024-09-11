package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// PluginMeta represents the metadata of a WordPress plugin
type PluginMeta struct {
	URL         string `default:"N/A"`
	Name        string `default:"Unknown"`
	Version     string `default:"0.0.0"`
	LastUpdated string `default:"N/A"`
	Installs    string `default:"N/A"`
	WPVersion   string `default:"N/A"`
	TestedUpTo  string `default:"N/A"`
	PHPVersion  string `default:"N/A"`
	Languages   string `default:"N/A"`
	Tags        string `default:"N/A"`
}

func main() {
	// Reset log file
	logFile, err := os.Create("scraper.log")
	if err != nil {
		log.Fatal("Failed to create log file:", err)
	}
	defer logFile.Close()

	// Set log output destination
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	log.Println("Starting scraping process")

	// Read CSV file containing URL list
	urls, err := readURLsFromCSV("plugin_urls.csv")
	if err != nil {
		log.Fatal("Failed to read URLs:", err)
	}

	log.Printf("Loaded %d URLs", len(urls))

	// Fetch plugin information for each URL
	var pluginMetas []PluginMeta
	for _, url := range urls {
		log.Printf("Processing URL: %s", url)
		meta, err := scrapePluginMetaWithRetry(url, 3) // Maximum 3 retries
		if err != nil {
			log.Printf("Warning: Error processing %s: %v", url, err)
		}
		pluginMetas = append(pluginMetas, meta)
		log.Printf("Completed processing URL: %s", url)
		time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second) // Random wait time of 1-5 seconds
	}

	// Export results to CSV
	err = exportToCSV(pluginMetas, "plugin_meta_results.csv")
	if err != nil {
		log.Fatal("Failed to export to CSV:", err)
	}

	log.Println("Scraping process completed")
	fmt.Println("Plugin metadata exported to CSV. Please check the log file for details.")
}

// scrapePluginMetaWithRetry attempts to scrape plugin metadata with retry logic
func scrapePluginMetaWithRetry(url string, maxRetries int) (PluginMeta, error) {
	var meta PluginMeta
	var err error

	for i := 0; i < maxRetries; i++ {
		meta, err = scrapePluginMeta(url)
		if err == nil {
			return meta, nil
		}

		if strings.Contains(err.Error(), "429") {
			retryAfter := time.Duration(30+rand.Intn(30)) * time.Second
			log.Printf("429 error. Retrying after %v: %s", retryAfter, url)
			time.Sleep(retryAfter)
		} else {
			return meta, err
		}
	}

	return meta, fmt.Errorf("maximum retry count reached: %v", err)
}

// scrapePluginMeta scrapes metadata from a single plugin page
func scrapePluginMeta(url string) (PluginMeta, error) {
	log.Printf("Starting scrape: %s", url)
	start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("HTTP GET request failed: %s", err)
		return PluginMeta{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Invalid HTTP status: %d for %s", resp.StatusCode, url)
		return PluginMeta{URL: url}, fmt.Errorf("invalid HTTP status: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Failed to parse HTML: %s", err)
		return PluginMeta{}, err
	}

	meta := PluginMeta{URL: url}
	meta.Name = strings.TrimSpace(doc.Find("h1.plugin-title").Text())

	doc.Find("div.entry-meta > div.widget.plugin-meta > ul > li").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		switch {
		case strings.Contains(text, "Version"):
			meta.Version = extractStrong(s)
		case strings.Contains(text, "Last updated"):
			meta.LastUpdated = extractStrong(s)
		case strings.Contains(text, "Active installations"):
			meta.Installs = extractStrong(s)
		case strings.Contains(text, "WordPress version"):
			meta.WPVersion = extractStrong(s)
		case strings.Contains(text, "Tested up to"):
			meta.TestedUpTo = extractStrong(s)
		case strings.Contains(text, "PHP version"):
			meta.PHPVersion = extractStrong(s)
		case strings.Contains(text, "Languages"):
			meta.Languages = strings.TrimSpace(s.Find("button").Text())
		case strings.Contains(text, "Tags"):
			meta.Tags = strings.TrimSpace(s.Find(".tags").Text())
		}
	})

	setDefaultValues(&meta)

	log.Printf("Completed scrape: %s (duration: %v)", url, time.Since(start))

	return meta, nil
}

// setDefaultValues sets default values for empty fields in PluginMeta
func setDefaultValues(meta *PluginMeta) {
	t := reflect.TypeOf(*meta)
	v := reflect.ValueOf(meta).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if defaultVal, ok := field.Tag.Lookup("default"); ok {
			if v.Field(i).String() == "" {
				v.Field(i).SetString(defaultVal)
				log.Printf("Set default value: %s=%s", field.Name, defaultVal)
			}
		}
	}
}

// extractStrong extracts the text within a strong element
func extractStrong(s *goquery.Selection) string {
	return strings.TrimSpace(s.Find("strong").Text())
}

// readURLsFromCSV reads plugin URLs from a CSV file
func readURLsFromCSV(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	
	// ヘッダー行を読み飛ばす
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	var urls []string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(record) > 0 {
			urls = append(urls, record[0])
		}
	}
	return urls, nil
}

// exportToCSV exports the scraped plugin metadata to a CSV file
func exportToCSV(data []PluginMeta, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"URL", "Name", "Version", "Last Updated", "Active Installations", "WordPress Version", "Tested Up To", "PHP Version", "Languages", "Tags"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, item := range data {
		row := []string{
			item.URL,
			item.Name,
			item.Version,
			item.LastUpdated,
			item.Installs,
			item.WPVersion,
			item.TestedUpTo,
			item.PHPVersion,
			item.Languages,
			item.Tags,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
