package wikipedia

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// StreamEdits returns all edits made to the german wiki for that filterFunc returns true.
// filterFunc gets the title of the article to make that decision
func StreamEdits(filterFunc func(event *Event) bool) <-chan Event {
	// Buffer of 25 should be more than enough
	var resultChannel = make(chan Event, 25)

	go func() {
		var (
			lastErrorTime time.Time
			backoff       int
		)

		for {
			log.Println("[StreamEdits] Connecting...")

			err := populateStreamEdits(filterFunc, resultChannel, func() {
				log.Println("[StreamEdits] Connected, processing events")
			})
			if err != nil {
				log.Printf("[StreamEdits] %s\n", err.Error())
			}

			if time.Since(lastErrorTime) < 5*time.Minute {
				backoff *= backoff
			} else {
				backoff = 2
			}

			lastErrorTime = time.Now()

			// Set wait time depending on how many fails there were, but reconnect within 5 minutes
			waitTime := time.Duration(backoff) * time.Second
			if waitTime > 5*time.Minute {
				waitTime = 5 * time.Minute
			}

			log.Printf("[StreamEdits] Waiting %s before reconnect...\n", waitTime)
			time.Sleep(waitTime)
		}
	}()

	return resultChannel
}

const (
	recentChangesURL = "https://stream.wikimedia.org/v2/stream/recentchange"
)

var noTimeoutClient = http.Client{}

func populateStreamEdits(filterFunc func(event *Event) bool, events chan<- Event, onConnect func()) (err error) {
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

		// If it's not an edit in the german wiki, we skip it.
		// Also skip bot edits and articles without titles (if they even exist?)
		if event.Bot || event.Type != "edit" || event.Wiki != "dewiki" || event.Title == "" {
			continue
		}

		if filterFunc(event) {
			events <- *event
		}
	}

	return
}
