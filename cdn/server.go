package cdn

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// EdgeServer represents a CDN edge server
type EdgeServer struct {
	Name         string
	Address      string
	CacheDir     string
	OriginServer string
	cache        map[string]*CacheEntry
	mu           sync.RWMutex
	server       *http.Server
}

// CacheEntry represents a cached file
type CacheEntry struct {
	FilePath   string
	Size       int64
	ETag       string
	LastAccess time.Time
	HitCount   int64
}

// NewEdgeServer creates a new CDN edge server
func NewEdgeServer(name, address, cacheDir, originServer string) *EdgeServer {
	return &EdgeServer{
		Name:         name,
		Address:      address,
		CacheDir:     cacheDir,
		OriginServer: originServer,
		cache:        make(map[string]*CacheEntry),
	}
}

// Start starts the edge server
func (e *EdgeServer) Start() error {
	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(e.CacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", e.handleRequest)
	mux.HandleFunc("/_stats", e.handleStats)
	mux.HandleFunc("/_cache/clear", e.handleCacheClear)

	e.server = &http.Server{
		Addr:    e.Address,
		Handler: mux,
	}

	log.Printf("CDN Edge Server '%s' starting on %s\n", e.Name, e.Address)
	log.Printf("Cache directory: %s\n", e.CacheDir)
	log.Printf("Origin server: %s\n", e.OriginServer)

	return e.server.ListenAndServe()
}

// Stop stops the edge server
func (e *EdgeServer) Stop() error {
	if e.server != nil {
		return e.server.Close()
	}
	return nil
}

func (e *EdgeServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	log.Printf("[%s] Request: %s\n", e.Name, path)

	// Check cache first
	e.mu.RLock()
	entry, exists := e.cache[path]
	e.mu.RUnlock()

	if exists {
		// Cache hit
		log.Printf("[%s] Cache HIT: %s\n", e.Name, path)
		e.serveCachedFile(w, entry)

		// Update stats
		e.mu.Lock()
		entry.LastAccess = time.Now()
		entry.HitCount++
		e.mu.Unlock()
		return
	}

	// Cache miss - fetch from origin
	log.Printf("[%s] Cache MISS: %s - fetching from origin\n", e.Name, path)
	e.fetchAndCache(w, r, path)
}

func (e *EdgeServer) serveCachedFile(w http.ResponseWriter, entry *CacheEntry) {
	file, err := os.Open(entry.FilePath)
	if err != nil {
		http.Error(w, "Error reading cached file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("X-Cache", "HIT")
	w.Header().Set("ETag", entry.ETag)
	w.Header().Set("Content-Type", getContentType(entry.FilePath))

	io.Copy(w, file)
}

func (e *EdgeServer) fetchAndCache(w http.ResponseWriter, r *http.Request, path string) {
	// Construct origin URL
	originURL := e.OriginServer + path

	// Fetch from origin server
	resp, err := http.Get(originURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching from origin: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Origin returned status: %d", resp.StatusCode), resp.StatusCode)
		return
	}

	// Create cache file path
	cacheFilePath := filepath.Join(e.CacheDir, generateCacheFileName(path))

	// Create cache file
	cacheFile, err := os.Create(cacheFilePath)
	if err != nil {
		log.Printf("Error creating cache file: %v\n", err)
		// Still serve the content even if caching fails
		w.Header().Set("X-Cache", "MISS")
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		io.Copy(w, resp.Body)
		return
	}
	defer cacheFile.Close()

	// Write to both cache and response
	w.Header().Set("X-Cache", "MISS")
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

	// Use TeeReader to write to both cache and response
	tee := io.TeeReader(resp.Body, cacheFile)
	written, _ := io.Copy(w, tee)

	// Generate ETag
	etag := generateETag(path)

	// Store in cache map
	e.mu.Lock()
	e.cache[path] = &CacheEntry{
		FilePath:   cacheFilePath,
		Size:       written,
		ETag:       etag,
		LastAccess: time.Now(),
		HitCount:   0,
	}
	e.mu.Unlock()

	log.Printf("[%s] Cached: %s (%d bytes)\n", e.Name, path, written)
}

func (e *EdgeServer) handleStats(w http.ResponseWriter, r *http.Request) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{\n")
	fmt.Fprintf(w, "  \"server_name\": \"%s\",\n", e.Name)
	fmt.Fprintf(w, "  \"address\": \"%s\",\n", e.Address)
	fmt.Fprintf(w, "  \"cached_items\": %d,\n", len(e.cache))
	fmt.Fprintf(w, "  \"cache\": [\n")

	i := 0
	for path, entry := range e.cache {
		if i > 0 {
			fmt.Fprintf(w, ",\n")
		}
		fmt.Fprintf(w, "    {\"path\": \"%s\", \"size\": %d, \"hits\": %d, \"last_access\": \"%s\"}",
			path, entry.Size, entry.HitCount, entry.LastAccess.Format(time.RFC3339))
		i++
	}

	fmt.Fprintf(w, "\n  ]\n")
	fmt.Fprintf(w, "}\n")
}

func (e *EdgeServer) handleCacheClear(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	// Delete cache files
	for _, entry := range e.cache {
		os.Remove(entry.FilePath)
	}

	// Clear cache map
	e.cache = make(map[string]*CacheEntry)

	log.Printf("[%s] Cache cleared\n", e.Name)
	fmt.Fprintf(w, "Cache cleared successfully\n")
}

func generateCacheFileName(path string) string {
	hash := md5.Sum([]byte(path))
	return fmt.Sprintf("%x", hash)
}

func generateETag(path string) string {
	hash := md5.Sum([]byte(path + time.Now().String()))
	return fmt.Sprintf("\"%x\"", hash)
}

func getContentType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	default:
		return "application/octet-stream"
	}
}
