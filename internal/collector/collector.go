package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type Collector struct {
	dataStoreBaseUrl string
	dataStorePath    string
	keptnApiToken    string
	keptnContext     string
	httpClient       *http.Client
}

type CollectorIface interface {
	GetEventsOfType(eventType string) ([]cloudevents.Event, error)
	GetTestStartedEvents() ([]cloudevents.Event, error)
	GetTestFinishedEvents() ([]cloudevents.Event, error)
	GetEvents() ([]cloudevents.Event, error)
	ParseEventsOfType(events []cloudevents.Event, filterType string) []cloudevents.Event
	MustParseEventsOfType(events []cloudevents.Event, filterType string) ([]cloudevents.Event, error)
	CollectExecutionIds(events []cloudevents.Event) ([]string, error)
	CollectBatchIds(events []cloudevents.Event) ([]string, error)
	CollectEarliestTime(events []cloudevents.Event) (time.Time, error)
	CollectLatestTime(events []cloudevents.Event) (time.Time, error)
}

type CollectedEvents struct {
	Events []cloudevents.Event `json:"events"`
}

type SyntheticTestFinishedEventData struct {
	SyntheticExecution struct {
		BatchId      string   `json:"batchId"`
		ExecutionIds []string `json:"executionIds"`
	} `json:"syntheticExecution"`
	Project string `json:"project"`
	Service string `json:"service"`
	Stage   string `json:"stage"`
}

func (c Collector) GetEventsOfType(eventType string) ([]cloudevents.Event, error) {
	u, err := url.Parse(c.dataStoreBaseUrl)
	if err != nil {
		return []cloudevents.Event{}, err
	}

	u.Path = c.dataStorePath
	query := u.Query()
	query.Add("keptnContext", c.keptnContext)

	if eventType != "" {
		query.Add("type", eventType)
	}

	u.RawQuery = query.Encode()

	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Set("x-token", c.keptnApiToken)
	req.Header.Set("accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return []cloudevents.Event{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	responseBody := CollectedEvents{}

	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return []cloudevents.Event{}, err
	}

	return responseBody.Events, nil
}

func (c Collector) GetTestStartedEvents() ([]cloudevents.Event, error) {
	return c.GetEventsOfType("sh.keptn.event.test.started")
}

func (c Collector) GetTestFinishedEvents() ([]cloudevents.Event, error) {
	return c.GetEventsOfType("sh.keptn.event.test.finished")
}

func (c Collector) GetEvents() ([]cloudevents.Event, error) {
	return c.GetEventsOfType("")
}

func (c Collector) ParseEventsOfType(events []cloudevents.Event, filterType string) []cloudevents.Event {
	eventsOfType := []cloudevents.Event{}

	for _, event := range events {
		if event.Type() == filterType {
			eventsOfType = append(eventsOfType, event)
		}
	}

	return eventsOfType
}

func (c Collector) MustParseEventsOfType(events []cloudevents.Event, filterType string) ([]cloudevents.Event, error) {
	eventsOfType := c.ParseEventsOfType(events, filterType)

	if len(events) < 1 {
		return eventsOfType, fmt.Errorf("no events found")
	}

	return eventsOfType, nil
}

func (c Collector) CollectExecutionIds(events []cloudevents.Event) ([]string, error) {
	executionIds := []string{}

	for _, event := range events {
		eventData := SyntheticTestFinishedEventData{}
		err := event.DataAs(&eventData)
		if err != nil {
			return []string{}, err
		}

		executionIds = append(executionIds, eventData.SyntheticExecution.ExecutionIds...)
	}

	return executionIds, nil
}

func (c Collector) CollectBatchIds(events []cloudevents.Event) ([]string, error) {
	batchIds := []string{}

	for _, event := range events {
		eventData := SyntheticTestFinishedEventData{}
		err := event.DataAs(&eventData)
		if err != nil {
			return []string{}, err
		}

		if eventData.SyntheticExecution.BatchId != "" {
			batchIds = append(batchIds, eventData.SyntheticExecution.BatchId)
		}
	}

	return batchIds, nil
}

func (c Collector) CollectEarliestTime(events []cloudevents.Event) (time.Time, error) {
	earliestTime := time.Time{}

	for _, event := range events {
		eventTime := event.Context.GetTime()

		if earliestTime.Equal(time.Time{}) || eventTime.Before(earliestTime) {
			earliestTime = eventTime
		}
	}

	if earliestTime.Equal(time.Time{}) {
		return earliestTime, fmt.Errorf("no timestamps found")
	}

	return earliestTime, nil
}

func (c Collector) CollectLatestTime(events []cloudevents.Event) (time.Time, error) {
	latestTime := time.Time{}

	for _, event := range events {
		eventTime := event.Context.GetTime()

		if latestTime.Equal(time.Time{}) || eventTime.After(latestTime) {
			latestTime = eventTime
		}
	}

	if latestTime.Equal(time.Time{}) {
		return latestTime, fmt.Errorf("no timestamps found")
	}

	return latestTime, nil
}

func NewCollector(keptnContext string) CollectorIface {
	dataStoreServiceHost := os.Getenv("MONGODB_DATASTORE_SERVICE_HOST")
	dataStoreServicePort := os.Getenv("MONGODB_DATASTORE_SERVICE_PORT")

	dataStoreScheme := os.Getenv("MONGODB_DATASTORE_SERVICE_SCHEME")
	if dataStoreScheme == "" {
		dataStoreScheme = "http"
	}

	dataStoreBaseUrl := fmt.Sprintf("%s://%s:%s", dataStoreScheme, dataStoreServiceHost, dataStoreServicePort)

	dataStorePath := os.Getenv("MONGODB_DATASTORE_PATH")
	if dataStorePath == "" {
		dataStorePath = "/event"
	}

	keptnApiToken := os.Getenv("KEPTN_API_TOKEN")
	httpClient := &http.Client{}

	return Collector{
		dataStoreBaseUrl,
		dataStorePath,
		keptnApiToken,
		keptnContext,
		httpClient,
	}
}
