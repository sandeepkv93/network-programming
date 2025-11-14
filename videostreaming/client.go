package videostreaming

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// StreamClient represents a video streaming client
type StreamClient struct {
	ServerURL string
}

// NewClient creates a new streaming client
func NewClient(serverURL string) *StreamClient {
	return &StreamClient{
		ServerURL: serverURL,
	}
}

// WatchStream connects to and watches a stream
func (c *StreamClient) WatchStream(streamID string, outputFile string) error {
	url := fmt.Sprintf("%s/stream/%s", c.ServerURL, streamID)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect to stream: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("stream not found or error: %d", resp.StatusCode)
	}

	fmt.Printf("Connected to stream: %s\n", streamID)
	fmt.Printf("Content-Type: %s\n", resp.Header.Get("Content-Type"))

	// Save stream to file or process
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %v", err)
		}
		defer file.Close()

		written, err := io.Copy(file, resp.Body)
		if err != nil {
			return fmt.Errorf("error receiving stream: %v", err)
		}

		fmt.Printf("Stream saved to %s (%d bytes)\n", outputFile, written)
	} else {
		// Just consume the stream
		buffer := make([]byte, 4096)
		total := 0
		for {
			n, err := resp.Body.Read(buffer)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			total += n
			fmt.Printf("\rReceived: %d bytes", total)
		}
		fmt.Println()
	}

	return nil
}

// ListStreams lists available streams
func (c *StreamClient) ListStreams() error {
	url := fmt.Sprintf("%s/api/streams", c.ServerURL)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get streams: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	fmt.Println("Available streams:")
	fmt.Println(string(body))
	return nil
}
