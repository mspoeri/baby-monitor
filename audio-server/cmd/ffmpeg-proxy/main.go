package main

import (
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {

	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Streaming audio")
		defer func() {
			slog.Info("Streaming audio finished")
		}()

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")

		cmd := exec.Command("ffmpeg", "-ar", "44100", "-ac", "1", "-f", "alsa", "-i", "plughw:1,0", "-f", "wav", "pipe:1")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			slog.Error("StdoutPipe error:", err)
			return
		}

		w.Header().Set("Content-Type", "audio/wav")

		err = cmd.Start()
		if err != nil {
			slog.Error("ffmpeg error:", err)
			err = cmd.Process.Signal(syscall.SIGTERM)
			if err != nil {
				slog.Error("Signal error:", err)
			}
			return
		}

		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			err := cmd.Process.Signal(syscall.SIGTERM)
			if err != nil {
				slog.Error("Signal error:", err)
			}
		}()

		go func() {
			defer func(stdout io.ReadCloser) {
				err := stdout.Close()
				if err != nil {
					slog.Error("Close error:", err)
				}
			}(stdout)
			n, err := io.Copy(w, stdout)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					slog.Error("io.Copy timeout error:", err)
				} else {
					slog.Error("io.Copy error:", err)
				}
				return
			}
			slog.Info("io.Copy wrote", n, "bytes")
		}()

		err = cmd.Wait()
		if err != nil {
			slog.Error("ffmpeg exited with error:", err)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		slog.Info("Not Found:", r.URL.Path)
	})

	slog.Info("Server is listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("ListenAndServe error:", err)
		os.Exit(1)
	}
}
