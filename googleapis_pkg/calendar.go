package googleapispkg

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func createEvent(service *calendar.Service, event *calendar.Event) (*calendar.Event, error) {
	createdEvent, err := service.Events.Insert("primary", event).Do()
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to create event: %w", err)
	}

	return createdEvent, nil
}

func initCalendarService(client *http.Client) (*calendar.Service, error) {
	ctx := context.Background()

	cal, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return cal, nil
}
