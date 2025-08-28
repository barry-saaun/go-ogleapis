package googleapispkg

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func CreateEvent(service *calendar.Service, title string, due time.Time) (*calendar.Event, error) {
	event := &calendar.Event{
		Summary: title,
		Start: &calendar.EventDateTime{
			DateTime: due.Format(time.RFC3339),
			TimeZone: "Australia/Melbourne",
		},

		End: &calendar.EventDateTime{
			DateTime: due.Format(time.RFC3339),
			TimeZone: "Australia/Melbourne",
		},
	}

	return service.Events.Insert("primary", event).Do()
}

func InitCalendarService(client *http.Client) (*calendar.Service, error) {
	ctx := context.Background()

	cal, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return cal, nil
}
