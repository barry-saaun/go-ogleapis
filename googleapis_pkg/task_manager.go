package googleapispkg

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/tasks/v1"
)

type AppEventMetadata struct {
	EventId    string `json:"eventId"`
	CalendarId string `json:"calendarId"`
}

const (
	metadataTag = "[APP_METADATA]"
	// Regex to find our metadata JSON in the notes
	metadataRegex = metadataTag + `\s*(\{.*\})`
	// Default event duration when creating a calendar event for a task
	defaultEventDuration = 30 * time.Minute
	// Timezone for events. Adjust if your users are in different timezones or allow customization.
	eventTimeZone = "Australia/Melbourne"

	defaultCalendarId = "primary"
)

type TaskManager struct {
	TaskTitle string
	Note      *string
	Completed *string
	DueTime   time.Time
}

type TaskManagerServices struct {
	tasksService    *tasks.Service
	calendarService *calendar.Service
}

type TaskResult struct {
	CreatedTask *tasks.Task
	TaskId      string
	EventId     string
	CalendarId  string
}

func InitTaskAndCalendarService(client *http.Client) (*TaskManagerServices, error) {
	calendarSrv, calErr := initCalendarService(client)
	tasksSrv, taskErr := initTasksService(client)

	if calErr != nil || taskErr != nil {
		return nil, nil
	}

	return &TaskManagerServices{
		tasksService:    tasksSrv,
		calendarService: calendarSrv,
	}, nil
}

func CreateTaskWithDueTime(ctx context.Context, taskManager *TaskManager, managerServices *TaskManagerServices, client *http.Client) (*TaskResult, error) {
	taskListId, err := GetTaskListId(managerServices.tasksService)
	if err != nil {
		return nil, err
	}

	event := &calendar.Event{
		Summary: taskManager.TaskTitle,
		Start: &calendar.EventDateTime{
			DateTime: taskManager.DueTime.Format(time.RFC3339),
			TimeZone: eventTimeZone,
		},
		End: &calendar.EventDateTime{
			DateTime: taskManager.DueTime.Format(time.RFC3339),
			TimeZone: eventTimeZone,
		},
	}

	createdEventRes, err := createEvent(managerServices.calendarService, event)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Created Calendar Event: %s (ID: %s, Link: %s)\n", createdEventRes.Summary, createdEventRes.Id, createdEventRes.HtmlLink)

	// --- Preparing metadata

	metadata := &AppEventMetadata{
		EventId:    createdEventRes.Id,
		CalendarId: defaultCalendarId,
	}

	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		// Clean up the created event if we fail to marshal metadata
		_ = managerServices.calendarService.Events.Delete(defaultCalendarId, createdEventRes.Id).Context(ctx).Do()
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	metadataString := fmt.Sprintf("%s %s", metadataTag, string(metadataJSON))

	combinedNotes := ""
	if taskManager.Note != nil {
		combinedNotes = *taskManager.Note + "\n"
	}
	combinedNotes += fmt.Sprintf("Related Event: %s\n%s", createdEventRes.HtmlLink, metadataString)

	// -----
	// fmt.Printf("combined notes: %s", combinedNotes)

	task := &tasks.Task{
		Title: taskManager.TaskTitle,
		Notes: combinedNotes,
		// Due:       taskManager.DueTime.Format("2006-01-02"),
		Due: taskManager.DueTime.Format(time.RFC3339),
		// Completed: nil,
	}

	createdTaskRes, err := createTask(managerServices.tasksService, task, taskListId)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Created Task ID: %s\n", createdTaskRes.Id)

	result := &TaskResult{
		CreatedTask: createdTaskRes,
		TaskId:      createdTaskRes.Id,
		EventId:     createdEventRes.Id,
	}

	return result, nil
}

func ModifyTask() {}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
