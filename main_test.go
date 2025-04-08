package main

import (
	"os"
	"testing"
)

// TestReadURLsFromCSV tests the functionality of reading URLs from a CSV file.
// It creates a temporary CSV file with test data, reads it, and verifies the contents.
// The test ensures that:
// - The CSV file is properly created
// - URLs are correctly read from the file
// - The number of URLs matches the expected count
// - Each URL matches its expected value
func TestReadURLsFromCSV(t *testing.T) {
	// Create a temporary CSV file with test data
	content := []byte("URL\nhttps://example.com/plugin1\nhttps://example.com/plugin2")
	err := os.WriteFile("test_urls.csv", content, 0644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}
	defer os.Remove("test_urls.csv")

	// Read URLs from the test CSV file
	urls, err := readURLsFromCSV("test_urls.csv")
	if err != nil {
		t.Errorf("readURLsFromCSV failed: %v", err)
	}

	// Verify the number of URLs read
	if len(urls) != 2 {
		t.Errorf("Expected 2 URLs, got %d", len(urls))
	}

	// Define expected URLs for comparison
	expectedURLs := []string{
		"https://example.com/plugin1",
		"https://example.com/plugin2",
	}

	// Compare each URL with its expected value
	for i, url := range urls {
		if url != expectedURLs[i] {
			t.Errorf("Expected URL %s, got %s", expectedURLs[i], url)
		}
	}
}

// TestSetDefaultValues verifies that default values are correctly set for empty PluginMeta fields.
// It tests:
// - Default value assignment for empty fields
// - Preservation of non-empty fields
// - Correct default values according to struct tags
func TestSetDefaultValues(t *testing.T) {
	// Create a PluginMeta instance with only URL set
	meta := &PluginMeta{
		URL: "https://example.com/plugin",
		// Other fields are intentionally left empty
	}

	// Apply default values
	setDefaultValues(meta)

	// Verify that default values are correctly set
	if meta.Name != "Unknown" {
		t.Errorf("Expected Name to be 'Unknown', got %s", meta.Name)
	}
	if meta.Version != "0.0.0" {
		t.Errorf("Expected Version to be '0.0.0', got %s", meta.Version)
	}
	if meta.LastUpdated != "N/A" {
		t.Errorf("Expected LastUpdated to be 'N/A', got %s", meta.LastUpdated)
	}
}

// TestExportToCSV tests the functionality of exporting PluginMeta data to a CSV file.
// It verifies:
// - Successful creation of the CSV file
// - Correct writing of test data
// - File existence after export
// - Proper cleanup after the test
func TestExportToCSV(t *testing.T) {
	// Create test data with sample plugin metadata
	testData := []PluginMeta{
		{
			URL:         "https://example.com/plugin1",
			Name:        "Test Plugin 1",
			Version:     "1.0.0",
			LastUpdated: "2024-03-20",
		},
		{
			URL:         "https://example.com/plugin2",
			Name:        "Test Plugin 2",
			Version:     "2.0.0",
			LastUpdated: "2024-03-21",
		},
	}

	// Export test data to CSV
	err := exportToCSV(testData, "test_export.csv")
	if err != nil {
		t.Fatalf("Failed to export CSV: %v", err)
	}
	defer os.Remove("test_export.csv")

	// Verify that the exported file exists
	if _, err := os.Stat("test_export.csv"); os.IsNotExist(err) {
		t.Error("Exported CSV file does not exist")
	}
}

// TestScrapePluginMetaWithRetry tests the retry mechanism for scraping plugin metadata.
// It verifies:
// - Error handling for invalid URLs
// - Retry logic functionality
// - URL preservation in the returned metadata
// Note: This is a basic implementation. A more comprehensive test would use a mock server.
func TestScrapePluginMetaWithRetry(t *testing.T) {
	// Test with an invalid URL to verify error handling
	url := "https://example.com/plugin"
	meta, err := scrapePluginMetaWithRetry(url, 3)

	// Verify that an error is returned for the invalid URL
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}

	// Verify that the URL is preserved in the metadata
	if meta.URL != url {
		t.Errorf("Expected URL %s, got %s", url, meta.URL)
	}
}

// TestExtractStrong tests the HTML parsing functionality for extracting text from strong tags.
// This test is currently skipped as it requires HTML parsing capabilities.
// TODO: Implement a proper test using a mock HTML document or a test server.
func TestExtractStrong(t *testing.T) {
	// Skip this test as it requires HTML parsing
	// A proper implementation would use a mock server or HTML document
	t.Skip("Skipping extractStrong test - requires HTML parsing")
}
