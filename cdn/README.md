## Content Delivery Network (CDN)

A Content Delivery Network is a geographically distributed network of servers that work together to provide fast delivery of Internet content. CDNs cache content at edge locations closer to users to reduce latency and improve performance.

## Table of Contents

1. [What is a CDN?](#what-is-a-cdn)
2. [How Does a CDN Work?](#how-does-a-cdn-work)
3. [Understanding the Code](#understanding-the-code)
4. [Further Reading](#further-reading)

### What is a CDN?

A Content Delivery Network (CDN) is a system of distributed servers that deliver web content to users based on their geographic location. The primary purpose is to reduce latency by serving content from a server closest to the user.

**Benefits**:
- **Reduced Latency**: Content served from nearby edge servers
- **Improved Performance**: Cached content loads faster
- **Reduced Bandwidth Costs**: Less traffic to origin server
- **Increased Reliability**: Redundancy across multiple servers
- **DDoS Protection**: Distributed architecture helps absorb attacks
- **Global Reach**: Serve users worldwide efficiently

**Common CDN Use Cases**:
- Static website content (images, CSS, JavaScript)
- Video streaming
- Software downloads
- API responses
- Dynamic content delivery

### How Does a CDN Work?

1. **Content Request**: User requests content from a website
2. **DNS Resolution**: DNS directs user to nearest edge server
3. **Cache Check**: Edge server checks if content is cached
4. **Cache Hit**: If cached, serve directly to user (fast)
5. **Cache Miss**: If not cached, fetch from origin server
6. **Cache and Serve**: Cache the content and serve to user
7. **Future Requests**: Subsequent requests are served from cache

**Architecture Components**:

```
User -> Edge Server (Cache) -> Origin Server
     <- (Cache Hit)
     <- (Cache Miss) ------> (Fetch)
                     <------ (Content)
     <- (Serve & Cache)
```

**Cache Strategy**:
- **Time-based**: Content expires after a certain time (TTL)
- **Event-based**: Content invalidated when updated
- **LRU (Least Recently Used)**: Remove least accessed items when cache is full

**Edge Server Locations**:
- Points of Presence (PoPs) distributed globally
- Each PoP contains multiple edge servers
- Strategic placement near major internet exchange points

### Understanding the Code

This implementation provides a simplified CDN with edge servers and an origin server.

#### Data Structures:

**EdgeServer**:
- `Name`: Identifier for the edge server
- `Address`: Address and port to listen on
- `CacheDir`: Directory to store cached files
- `OriginServer`: URL of the origin server
- `cache`: Map of cached entries

**CacheEntry**:
- `FilePath`: Local path to cached file
- `Size`: Size of cached file
- `ETag`: Entity tag for cache validation
- `LastAccess`: When the entry was last accessed
- `HitCount`: Number of cache hits

**OriginServer**:
- `Address`: Address and port to listen on
- Simple server that serves the original content

#### Functions:

**EdgeServer**:
- `NewEdgeServer(name, address, cacheDir, originServer string) *EdgeServer`: Creates edge server
- `Start() error`: Starts the edge server
- `Stop() error`: Stops the edge server
- `handleRequest(w, r)`: Handles incoming requests, checks cache, fetches if needed
- `serveCachedFile(w, entry)`: Serves content from cache
- `fetchAndCache(w, r, path)`: Fetches from origin and caches
- `handleStats(w, r)`: Provides cache statistics
- `handleCacheClear(w, r)`: Clears the cache

**OriginServer**:
- `NewOriginServer(address string) *OriginServer`: Creates origin server
- `Start() error`: Starts the origin server
- `handleRequest(w, r)`: Serves original content

#### Features:

- In-memory cache map with disk storage
- Cache hit/miss tracking
- Cache statistics endpoint (`/_stats`)
- Cache clearing endpoint (`/_cache/clear`)
- ETag generation for cache validation
- Content-Type detection
- Concurrent request handling
- Automatic cache directory creation

#### Example Usage:

```go
// Start origin server
origin := NewOriginServer(":8080")
go origin.Start()

// Start edge servers in different locations
edge1 := NewEdgeServer("US-East", ":9001", "/tmp/cdn-cache-1", "http://localhost:8080")
edge2 := NewEdgeServer("EU-West", ":9002", "/tmp/cdn-cache-2", "http://localhost:8080")

go edge1.Start()
go edge2.Start()

// Access content through edge servers:
// http://localhost:9001/page.html  (first request: cache miss)
// http://localhost:9001/page.html  (subsequent: cache hit)

// View statistics:
// http://localhost:9001/_stats

// Clear cache:
// POST http://localhost:9001/_cache/clear
```

### Further Reading

- [Content Delivery Network - Wikipedia](https://en.wikipedia.org/wiki/Content_delivery_network)
- [How CDNs Work](https://www.cloudflare.com/learning/cdn/what-is-a-cdn/)
- [HTTP Caching](https://developer.mozilla.org/en-US/docs/Web/HTTP/Caching)
- [Cache-Control Headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control)
- [CDN Architecture](https://www.akamai.com/our-thinking/cdn/what-is-a-cdn)
