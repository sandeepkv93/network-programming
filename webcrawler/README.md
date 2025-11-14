## Web Crawler

A web crawler (also known as a web spider) is a program that automatically browses the World Wide Web in a methodical, automated manner. It is commonly used for web indexing, data mining, and monitoring websites.

## Table of Contents

1. [What is a Web Crawler?](#what-is-a-web-crawler)
2. [How Does Web Crawling Work?](#how-does-web-crawling-work)
3. [Understanding the Code](#understanding-the-code)
4. [Further Reading](#further-reading)

### What is a Web Crawler?

A web crawler is an automated program that systematically browses the internet, typically for the purpose of web indexing. Search engines like Google use web crawlers to discover and index web pages. The crawler starts with a seed URL and follows hyperlinks to discover new pages.

**Key Components**:
- **URL Queue**: A queue of URLs to be crawled
- **Visited Set**: A set to track already visited URLs to avoid cycles
- **Parser**: Extracts links and content from HTML pages
- **Politeness Policy**: Rules to avoid overloading servers (rate limiting, respecting robots.txt)

### How Does Web Crawling Work?

1. **Initialize**: Start with a seed URL or list of URLs
2. **Fetch**: Download the HTML content of the URL
3. **Parse**: Extract links and relevant data from the HTML
4. **Filter**: Filter out already visited URLs and off-domain links
5. **Queue**: Add new URLs to the crawl queue
6. **Repeat**: Continue until the queue is empty or limits are reached

**Depth-First vs Breadth-First**:
- **Depth-First**: Follows links deep into the site structure before exploring siblings
- **Breadth-First**: Explores all links at the current depth before going deeper

**Considerations**:
- **Robots.txt**: Respect the website's crawling rules
- **Rate Limiting**: Avoid overwhelming the server with too many requests
- **Duplicate Detection**: Avoid crawling the same page multiple times
- **URL Normalization**: Handle different URL formats that point to the same resource

### Understanding the Code

#### Data Structures:

- `Crawler`: The main crawler structure with configuration options:
  - `MaxDepth`: Maximum depth to crawl from the starting URL
  - `MaxPages`: Maximum number of pages to crawl
  - `Timeout`: HTTP request timeout
  - `visited`: Map to track visited URLs
  - `OnPage`: Callback function when a page is successfully crawled
  - `OnError`: Callback function when an error occurs

#### Functions:

- `NewCrawler(maxDepth, maxPages int) *Crawler`: Creates a new crawler instance
- `Crawl(startURL string) error`: Starts the crawling process from the given URL
- `crawlRecursive(pageURL, baseDomain string, depth int)`: Recursively crawls pages
- `fetchPage(pageURL string) (string, error)`: Fetches the HTML content of a page
- `extractLinks(body, baseURL, baseDomain string) []string`: Extracts links from HTML
- `makeAbsolute(link, baseURL string) string`: Converts relative URLs to absolute URLs
- `GetVisitedURLs() []string`: Returns all visited URLs

#### Features:

- Configurable depth and page limits
- Concurrent crawling with semaphore for rate limiting
- Same-domain filtering (only crawls pages from the starting domain)
- Duplicate URL detection
- Callback support for custom processing
- User-agent customization

### Further Reading

- [Web Crawler - Wikipedia](https://en.wikipedia.org/wiki/Web_crawler)
- [How Search Engines Work](https://www.google.com/search/howsearchworks/crawling-indexing/)
- [Robots Exclusion Protocol](https://www.robotstxt.org/)
- [Politeness Policy for Web Crawlers](https://en.wikipedia.org/wiki/Web_crawler#Politeness_policy)
