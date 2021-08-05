package transfer

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Notifier defines the API for notifying a remote care organization that an eOverdracht FHIR task has been updated,
type Notifier interface {
	// Notify sends a notification to the given endpoint.
	Notify(endpoint, taskOwnerDID string) error
}

// fireAndForgetNotifier is a notifier that is optimistic about the receiver's availability.
// It just sends the notification and assumes the receiver is available.
type fireAndForgetNotifier struct {

}

func (f fireAndForgetNotifier) Notify(endpoint, taskOwnerDID string) error {
	client := http.Client{}
	notificationURL := fmt.Sprintf("%s?taskOwnerDID=%s", endpoint, url.QueryEscape(taskOwnerDID))
	response, err := client.Post(notificationURL, "", bytes.NewReader([]byte{}))
	if response != nil {
		// We try to be a nice client by always reading the HTTP response
		_, _ = io.Copy(io.Discard, response.Body)
	}
	if err != nil {
		return fmt.Errorf("unable to send eOverdracht notification (url=%s): %w", endpoint, err)
	}
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return fmt.Errorf("eOverdracht notification endpoint returned non-OK error code (url=%s,status-code=%d)", endpoint, response.StatusCode)
	}
	return nil
}

