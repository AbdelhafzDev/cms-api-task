package httputil

import (
	"context"
	"errors"
	"net/http"
)

var ErrStreamingNotSupported = errors.New("streaming not supported")

func BinaryStream(ctx context.Context, w http.ResponseWriter, contentType string, chunks <-chan []byte, errChan <-chan error) error {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Accept-Ranges", "bytes")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return ErrStreamingNotSupported
	}

	for {
		select {
		case chunk, ok := <-chunks:
			if !ok {
				return nil
			}
			if _, err := w.Write(chunk); err != nil {
				return err
			}
			flusher.Flush()

		case err := <-errChan:
			if err != nil {
				return err
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
