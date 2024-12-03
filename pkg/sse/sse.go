package sse

import (
	"encoding/json"
	"io"
	"net/http"
)

// WriteDone writes a [DONE] SSE message to the writer.
func WriteDone(w io.Writer) {
	_, _ = w.Write([]byte("data: [DONE]\n\n"))
}

// WriteData writes a data SSE message to the writer.
func WriteData(w io.Writer, v any) error {
	_, _ = w.Write([]byte("data: "))
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	_, _ = w.Write([]byte("\n")) // Encode() adds one newline, so add only one more here.
	return nil
}

// WriteDataAndFlush writes a data SSE message to the writer and flushes the writer.
func WriteDataAndFlush(w io.Writer, v any) error {
	if err := WriteData(w, v); err != nil {
		return err
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}

func WriteRawData(w io.Writer, raw string) error {
	_, _ = w.Write([]byte("data: " + raw + "\n\n"))
	return nil
}

// WriteEvent writes a data SSE message to the writer.
func WriteEvent(w io.Writer, name string) error {
	_, err := w.Write([]byte("event: " + name))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\n"))
	if err != nil {
		return err
	}
	return nil
}

// WriteEventAndFlush writes a data SSE message to the writer and flushes the writer.
func WriteEventAndFlush(w io.Writer, name string) error {
	if err := WriteEvent(w, name); err != nil {
		return err
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}

// WriteStreamingHeaders writes the headers for a streaming response.
func WriteStreamingHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")
}
