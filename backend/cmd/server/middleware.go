package main

import (
	"io"
	"log"
	"log/slog"
	"net/http"
	"time"
)

// LOGGING

// code modified from https://gist.github.com/arunvelsriram/74fa35d1c9e6bbf94b8f21b42a5f07e1
type (
	// struct for holding response details
	responseData struct {
		status int
		size   int
	}

	// our http.ResponseWriter implementation
	loggingResponseWriter struct {
		http.ResponseWriter // compose original http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b) // write response using original http.ResponseWriter
	r.responseData.size += size            // capture size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode) // write status code using original http.ResponseWriter
	r.responseData.status = statusCode       // capture status code
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		url := r.RequestURI
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		requestBody := string(b)

		method := r.Method

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lrw := loggingResponseWriter{
			ResponseWriter: w, // compose original http.ResponseWriter
			responseData:   responseData,
		}
		next.ServeHTTP(&lrw, r) // inject our implementation of http.ResponseWriter

		duration := time.Since(start)

		slog.Info("Handled request", "method", method, "url", url, "requestBody", requestBody, "requestHeaders", r.Header, "duration", duration, "statusCode", responseData.status, "responseSize", responseData.size)
	})
}
