package videostreaming

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// StreamServer represents a video streaming server
type StreamServer struct {
	Address string
	streams map[string]*Stream
	mu      sync.RWMutex
	server  *http.Server
}

// Stream represents a video stream
type Stream struct {
	ID          string
	Title       string
	Description string
	StartTime   time.Time
	Viewers     int
	mu          sync.Mutex
}

// NewStreamServer creates a new video streaming server
func NewStreamServer(address string) *StreamServer {
	return &StreamServer{
		Address: address,
		streams: make(map[string]*Stream),
	}
}

// Start starts the streaming server
func (s *StreamServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/stream/", s.handleStream)
	mux.HandleFunc("/api/streams", s.handleListStreams)
	mux.HandleFunc("/api/stream/create", s.handleCreateStream)
	mux.HandleFunc("/", s.handleIndex)

	s.server = &http.Server{
		Addr:    s.Address,
		Handler: mux,
	}

	log.Printf("Video Streaming Server started on %s\n", s.Address)
	return s.server.ListenAndServe()
}

// Stop stops the streaming server
func (s *StreamServer) Stop() error {
	if s.server != nil {
		return s.server.Close()
	}
	return nil
}

func (s *StreamServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Video Streaming Server</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .stream { background: #f0f0f0; padding: 15px; margin: 10px 0; border-radius: 5px; }
    </style>
</head>
<body>
    <h1>Video Streaming Server</h1>
    <p>Active streams are listed below.</p>
    <div id="streams"></div>
    <script>
        fetch('/api/streams')
            .then(r => r.json())
            .then(streams => {
                const div = document.getElementById('streams');
                streams.forEach(stream => {
                    div.innerHTML += '<div class="stream">' +
                        '<h3>' + stream.title + '</h3>' +
                        '<p>' + stream.description + '</p>' +
                        '<p>Viewers: ' + stream.viewers + '</p>' +
                        '<a href="/stream/' + stream.id + '">Watch Stream</a>' +
                        '</div>';
                });
            });
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func (s *StreamServer) handleStream(w http.ResponseWriter, r *http.Request) {
	streamID := r.URL.Path[len("/stream/"):]

	s.mu.RLock()
	stream, exists := s.streams[streamID]
	s.mu.RUnlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	// Increment viewer count
	stream.mu.Lock()
	stream.Viewers++
	viewers := stream.Viewers
	stream.mu.Unlock()

	log.Printf("New viewer for stream %s (total: %d)\n", streamID, viewers)

	// Set headers for streaming
	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Simulate streaming by sending data chunks
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	// Send video chunks (simulated)
	for i := 0; i < 100; i++ {
		// In a real implementation, read from video file or live source
		chunk := make([]byte, 4096)
		w.Write(chunk)
		flusher.Flush()
		time.Sleep(100 * time.Millisecond)

		// Check if client disconnected
		select {
		case <-r.Context().Done():
			stream.mu.Lock()
			stream.Viewers--
			stream.mu.Unlock()
			log.Printf("Viewer left stream %s\n", streamID)
			return
		default:
		}
	}

	stream.mu.Lock()
	stream.Viewers--
	stream.mu.Unlock()
}

func (s *StreamServer) handleListStreams(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "[")

	i := 0
	for _, stream := range s.streams {
		if i > 0 {
			fmt.Fprint(w, ",")
		}
		fmt.Fprintf(w, `{"id":"%s","title":"%s","description":"%s","viewers":%d}`,
			stream.ID, stream.Title, stream.Description, stream.Viewers)
		i++
	}

	fmt.Fprint(w, "]")
}

func (s *StreamServer) handleCreateStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	streamID := r.FormValue("id")
	title := r.FormValue("title")
	description := r.FormValue("description")

	if streamID == "" {
		http.Error(w, "Stream ID required", http.StatusBadRequest)
		return
	}

	stream := &Stream{
		ID:          streamID,
		Title:       title,
		Description: description,
		StartTime:   time.Now(),
		Viewers:     0,
	}

	s.mu.Lock()
	s.streams[streamID] = stream
	s.mu.Unlock()

	log.Printf("Created stream: %s - %s\n", streamID, title)
	fmt.Fprintf(w, "Stream created: %s\n", streamID)
}
