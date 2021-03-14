package wikipedia

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// StreamEdits returns all edits made to the german wiki for that filterFunc returns true.
// FilterFunc should return true if the article with the given title should be sent to events
func StreamEdits(filterFunc func(title string) bool) <-chan Event {
	var resultChannel = make(chan Event, 250)

	go func() {
		var (
			lastErrorTime time.Time
			backoff       int
		)

		for {
			log.Println("[StreamEdits]: Connecting...")

			err := populateStreamEdits(filterFunc, resultChannel, func() {
				log.Println("[StreamEdits]: Connected, processing events")
			})
			if err != nil {
				log.Printf("[StreamEdits]: %s\n", err.Error())
			}

			if time.Since(lastErrorTime) < 5*time.Minute {
				backoff *= backoff
			} else {
				backoff = 2
			}

			lastErrorTime = time.Now()

			waitTime := time.Duration(backoff) * time.Second
			log.Printf("[StreamEdits]: Waiting %s before reconnect...\n", waitTime)
			time.Sleep(waitTime)
		}
	}()

	return resultChannel
}

const (
	recentChangesURL = "https://stream.wikimedia.org/v2/stream/recentchange"
)

var noTimeoutClient = http.Client{}

func populateStreamEdits(filterFunc func(title string) bool, events chan<- Event, onConnect func()) (err error) {
	req, err := http.NewRequest(http.MethodGet, recentChangesURL, nil)
	if err != nil {
		return
	}

	req.Header.Set("Accept", "application/json")

	resp, err := noTimeoutClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, resp.Status)
	}

	onConnect()

	dec := json.NewDecoder(resp.Body)

	var event = new(Event)

	for dec.More() {
		err = dec.Decode(event)
		if err != nil {
			break
		}

		// If it's not an edit in the german wiki, we skip it
		if event.Bot || event.Type != "edit" || event.Wiki != "dewiki" || event.Title == "" {
			continue
		}

		if filterFunc(event.Title) {
			events <- *event
		}
	}

	return
}
