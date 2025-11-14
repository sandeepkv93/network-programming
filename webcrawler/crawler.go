package webcrawler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Crawler represents a web crawler
type Crawler struct {
	MaxDepth     int
	MaxPages     int
	Timeout      time.Duration
	UserAgent    string
	visited      map[string]bool
	mu           sync.Mutex
	pageCount    int
	OnPage       func(url string, depth int, body string)
	OnError      func(url string, err error)
}

// NewCrawler creates a new web crawler
func NewCrawler(maxDepth, maxPages int) *Crawler {
	return &Crawler{
		MaxDepth:  maxDepth,
		MaxPages:  maxPages,
		Timeout:   30 * time.Second,
		UserAgent: "Go-WebCrawler/1.0",
		visited:   make(map[string]bool),
	}
}

// Crawl starts crawling from the given URL
func (c *Crawler) Crawl(startURL string) error {
	// Validate URL
	u, err := url.Parse(startURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}

	fmt.Printf("Starting web crawler from %s\n", startURL)
	fmt.Printf("Max depth: %d, Max pages: %d\n", c.MaxDepth, c.MaxPages)

	c.crawlRecursive(startURL, u.Host, 0)

	fmt.Printf("\nCrawling completed. Total pages visited: %d\n", c.pageCount)
	return nil
}

func (c *Crawler) crawlRecursive(pageURL, baseDomain string, depth int) {
	// Check depth limit
	if depth > c.MaxDepth {
		return
	}

	// Check if we've reached max pages
	c.mu.Lock()
	if c.pageCount >= c.MaxPages {
		c.mu.Unlock()
		return
	}
	c.mu.Unlock()

	// Check if already visited
	c.mu.Lock()
	if c.visited[pageURL] {
		c.mu.Unlock()
		return
	}
	c.visited[pageURL] = true
	c.pageCount++
	currentCount := c.pageCount
	c.mu.Unlock()

	fmt.Printf("[%d/%d] Crawling (depth %d): %s\n", currentCount, c.MaxPages, depth, pageURL)

	// Fetch the page
	body, err := c.fetchPage(pageURL)
	if err != nil {
		if c.OnError != nil {
			c.OnError(pageURL, err)
		}
		return
	}

	// Call the page callback if set
	if c.OnPage != nil {
		c.OnPage(pageURL, depth, body)
	}

	// Extract links from the page
	links := c.extractLinks(body, pageURL, baseDomain)

	// Crawl each link
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 5) // Limit concurrent requests

	for _, link := range links {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			c.crawlRecursive(l, baseDomain, depth+1)
		}(link)
	}

	wg.Wait()
}

func (c *Crawler) fetchPage(pageURL string) (string, error) {
	client := &http.Client{
		Timeout: c.Timeout,
	}

	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), nil
}

func (c *Crawler) extractLinks(body, baseURL, baseDomain string) []string {
	var links []string
	seen := make(map[string]bool)

	// Regular expression to find links
	re := regexp.MustCompile(`href=["']([^"']+)["']`)
	matches := re.FindAllStringSubmatch(body, -1)

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		link := match[1]

		// Skip anchors and javascript
		if strings.HasPrefix(link, "#") || strings.HasPrefix(link, "javascript:") {
			continue
		}

		// Convert relative URLs to absolute
		absoluteURL := c.makeAbsolute(link, baseURL)
		if absoluteURL == "" {
			continue
		}

		// Only crawl links from the same domain
		u, err := url.Parse(absoluteURL)
		if err != nil || u.Host != baseDomain {
			continue
		}

		// Avoid duplicates
		if !seen[absoluteURL] {
			seen[absoluteURL] = true
			links = append(links, absoluteURL)
		}
	}

	return links
}

func (c *Crawler) makeAbsolute(link, baseURL string) string {
	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	u, err := url.Parse(link)
	if err != nil {
		return ""
	}

	return base.ResolveReference(u).String()
}

// GetVisitedURLs returns all visited URLs
func (c *Crawler) GetVisitedURLs() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	urls := make([]string, 0, len(c.visited))
	for url := range c.visited {
		urls = append(urls, url)
	}
	return urls
}
