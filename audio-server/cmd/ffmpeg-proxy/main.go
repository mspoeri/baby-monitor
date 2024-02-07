package main

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"sync"
)

var (
	ffmpegStdout io.ReadCloser
	ffmpegCmd    *exec.Cmd
	buffer       = make([]byte, 0)
	bufferMutex  = &sync.RWMutex{}
)

func main() {
	ffmpegCmd = exec.Command("ffmpeg", "-ar", "44100", "-ac", "1", "-f", "alsa", "-i", "plughw:1,0", "-f", "wav", "pipe:1")
	var err error
	ffmpegStdout, err = ffmpegCmd.StdoutPipe()
	if err != nil {
		slog.Error("StdoutPipe error:", err)
		os.Exit(1)
		return
	}

	err = ffmpegCmd.Start()
	if err != nil {
		slog.Error("ffmpeg error:", err)
		os.Exit(1)
		return
	}

	go func() {
		for {
			tmp := make([]byte, 1024)
			n, err := ffmpegStdout.Read(tmp)
			if err != nil {
				slog.Error("Read error:", err)
				os.Exit(1)
			}

			bufferMutex.Lock()
			buffer = append(buffer, tmp[:n]...)
			bufferMutex.Unlock()
		}
	}()

	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")

		go func() {
			localBuffer := make([]byte, 0)
			for {
				bufferMutex.RLock()
				if len(buffer) > len(localBuffer) {
					newBytes := buffer[len(localBuffer):]
					localBuffer = append(localBuffer, newBytes...)
				}
				bufferMutex.RUnlock()

				if len(localBuffer) > 0 {
					_, err := w.Write(localBuffer)
					if err != nil {
						slog.Error("Write error:", err)
						return
					}
					localBuffer = localBuffer[:0]
				}
			}
		}()
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		slog.Info("Not Found:", r.URL.Path)
	})

	slog.Info("Server is listening on port 8081")
	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		slog.Error("ListenAndServe error:", err)
		os.Exit(1)
	}
}
