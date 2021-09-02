package transfer

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/gommon/log"
)

// Notifier defines the API for notifying a remote care organization that an eOverdracht FHIR task has been updated,
type Notifier interface {
	// Notify sends a notification to the given endpoint.
	Notify(token, endpoint string) error
}

// fireAndForgetNotifier is a notifier that is optimistic about the receiver's availability.
// It just sends the notification and assumes the receiver is available.
type fireAndForgetNotifier struct {
}

func (f fireAndForgetNotifier) Notify(token, endpoint string) error {
	response, err := resty.New().
		R().
		SetBody([]byte{}).
		SetAuthToken(token).
		Post(endpoint)
	if err != nil {
		return fmt.Errorf("unable to send eOverdracht notification (url=%s): %w", endpoint, err)
	}

	if !response.IsSuccess() {
		log.Warnf("Server response: %s", response.String())
		return fmt.Errorf("eOverdracht notification endpoint returned non-OK error code (url=%s,status-code=%d)", endpoint, response.StatusCode())
	}

	return nil
}
