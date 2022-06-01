package collector

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"gotest.tools/assert"
)

func newMockTestFinishedEvent() cloudevents.Event {
	mockTestFinishedEvent := cloudevents.NewEvent()
	mockTestFinishedEvent.SetType("sh.keptn.event.test.finished")
	mockTestFinishedEvent.SetSpecVersion("1.0")

	mockTestFinishedEvent.DataEncoded = []byte(`{
		"syntheticExecution": {
			"batchId": "8602313944601341093",
			"executionIds": [
				"936772897",
				"936772898",
				"936772899"
			]
		},
		"project": "simplenode-gitlab",
		"service": "simplenodeservice",
		"stage": "staging"
	}`)

	return mockTestFinishedEvent
}

func newMockTestStartedEvent() cloudevents.Event {
	mockTestFinishedEvent := cloudevents.NewEvent()
	mockTestFinishedEvent.SetType("sh.keptn.event.test.started")
	mockTestFinishedEvent.SetSpecVersion("1.0")

	mockTestFinishedEvent.DataEncoded = []byte(`{
		"project": "simplenode-gitlab",
		"service": "simplenodeservice",
		"stage": "staging"
	}`)

	return mockTestFinishedEvent
}

func TestGetEvents(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/event" {
			t.Errorf("Expected to request '/event', got: %s", r.URL.Path)
		}

		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"events":[{ "specversion": "1.0" },{ "specversion": "1.0" }]}`))
	}))
	defer server.Close()

	c := Collector{
		dataStoreBaseUrl: server.URL,
		dataStorePath:    "/event",
		keptnApiToken:    "token",
		httpClient:       server.Client(),
	}

	events, err := c.GetEvents("keptnContext")
	assert.NilError(t, err)
	assert.Equal(t, len(events), 2)
}

func TestParseEventsOfType(t *testing.T) {
	c := NewCollector()

	events := c.ParseEventsOfType([]cloudevents.Event{}, "")
	assert.Equal(t, len(events), 0)
}

func TestMustParseEventsOfType(t *testing.T) {
	c := NewCollector()

	events, err := c.MustParseEventsOfType([]cloudevents.Event{}, "sh.keptn.event.test.finished")
	assert.Error(t, err, "no events found")
	assert.Equal(t, len(events), 0)

	events, err = c.MustParseEventsOfType([]cloudevents.Event{newMockTestFinishedEvent(), newMockTestStartedEvent()}, "sh.keptn.event.test.finished")
	assert.NilError(t, err)
	assert.Equal(t, len(events), 1)
}

func TestCollectExecutionIds(t *testing.T) {
	c := NewCollector()

	executionIds, err := c.CollectExecutionIds([]cloudevents.Event{newMockTestFinishedEvent(), newMockTestStartedEvent()})
	assert.NilError(t, err)
	assert.Equal(t, len(executionIds), 3)
}

func TestCollectBatchIds(t *testing.T) {
	c := NewCollector()

	batchIds, err := c.CollectBatchIds([]cloudevents.Event{newMockTestFinishedEvent(), newMockTestStartedEvent()})
	assert.NilError(t, err)
	assert.Equal(t, len(batchIds), 1)
}

func TestCollectEarliestTime(t *testing.T) {
	collector := NewCollector()

	a := cloudevents.NewEvent()
	timestampA, _ := time.Parse(time.RFC3339, "2022-04-07T12:04:28Z")
	a.SetTime(timestampA)

	b := cloudevents.NewEvent()
	timestampB, _ := time.Parse(time.RFC3339, "2022-04-07T12:05:28Z")
	b.SetTime(timestampB)

	c := cloudevents.NewEvent()
	timestampC, _ := time.Parse(time.RFC3339, "2022-04-07T12:06:28Z")
	c.SetTime(timestampC)

	earliestTimestamp, err := collector.CollectEarliestTime([]cloudevents.Event{b, a, c})
	assert.NilError(t, err)
	assert.Equal(t, timestampA, earliestTimestamp)

	latestTimestamp, err := collector.CollectLatestTime([]cloudevents.Event{b, a, c})
	assert.NilError(t, err)
	assert.Equal(t, timestampC, latestTimestamp)
}
