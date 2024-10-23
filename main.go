package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Watcher struct {
	dir      string
	previous map[string]time.Time
	clinets  []chan string
}

func NewWatcher(dir string) *Watcher {
	return &Watcher{
		dir:      dir,
		previous: make(map[string]time.Time),
		clinets:  make([]chan string, 0),
	}
}

func (w *Watcher) StartWatching() {
	for {
		err := w.checkForChanges()

		if err != nil {
			fmt.Println("Error chacking for changes: ", err)
		}

		time.Sleep(2 * time.Second)
	}
}

func (w *Watcher) checkForChanges() error {
	err := filepath.Walk(w.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			lastModified := info.ModTime()

			if prevTime, exists := w.previous[path]; !exists {
				w.notifyClients(fmt.Sprintf("File created: &s", path))
			} else if lastModified.After(prevTime) {
				w.notifyClients(fmt.Sprintf("File modified: &s", path))
			}

			w.previous[path] = lastModified
		}

		return nil
	})

	for path := range w.previous {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			w.notifyClients(fmt.Sprintf("File deleted: %s", path))
			delete(w.previous, path)
		}
	}

	return err
}

func (w *Watcher) notifyClients(message string) {
	for _, client := range w.clinets {
		client <- message
	}
}

func (w *Watcher) HandleSSEConnection(wr http.ResponseWriter, r *http.Request) {
	clientChan := make(chan string)
	w.clinets = append(w.clinets, clientChan)

	defer func() {
		close(clientChan)
		for i, c := range w.clinets {
			if c == clientChan {
				w.clinets = append(w.clinets[:i], w.clinets[i+1:]...)
				break
			}
		}
	}()

	wr.Header().Set("Content-Type", "text/event-stream")
	wr.Header().Set("Cache-Control", "no-cache")
	wr.Header().Set("Connection", "keep-alive")

	for msg := range clientChan {
		fmt.Fprintf(wr, "data: %s\n\n", msg)
		wr.(http.Flusher).Flush()
	}

}

func ServeHTML(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>File Watcher</title>
		<script>
			const eventSource = new EventSource("/events");
			eventSource.onmessage = function(event) {
				const message = document.createElement("div");
				message.textContent = event.data;
				document.body.appendChild(message);
			};
		</script>
	</head>
	<body>
		<h1>File Watcher</h1>
		<div id="messages"></div>
	</body>
	</html>
	`
	w.Write([]byte(html))
}

func main() {
	dir := "./watched"

	watcher := NewWatcher(dir)

	go watcher.StartWatching()

	http.HandleFunc("/", ServeHTML)
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		watcher.HandleSSEConnection(w, r)
	})

	fmt.Println("Server start at: 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

}
