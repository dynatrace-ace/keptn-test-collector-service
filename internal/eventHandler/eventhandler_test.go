//go:generate mockgen -source=../collector/collector.go -destination=collector_mock_test.go -package=eventHandler CollectorIface

package eventHandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	time "time"

	gomock "github.com/golang/mock/gomock"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0/fake"
	"gotest.tools/assert"

	keptn "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	cloudevents "github.com/cloudevents/sdk-go/v2" // make sure to use v2 cloudevents here
)

/**
 * loads a cloud event from the passed test json file and initializes a keptn object with it
 */
func initializeTestObjects(eventFileName string) (*keptnv2.Keptn, *cloudevents.Event, error) {
	// load sample event
	eventFile, err := ioutil.ReadFile(eventFileName)
	if err != nil {
		return nil, nil, fmt.Errorf("Cant load %s: %s", eventFileName, err.Error())
	}

	incomingEvent := &cloudevents.Event{}
	err = json.Unmarshal(eventFile, incomingEvent)
	if err != nil {
		return nil, nil, fmt.Errorf("Error parsing: %s", err.Error())
	}

	// Add a Fake EventSender to KeptnOptions
	var keptnOptions = keptn.KeptnOpts{
		EventSender: &fake.EventSender{},
	}
	keptnOptions.UseLocalFileSystem = true
	myKeptn, err := keptnv2.NewKeptn(incomingEvent, keptnOptions)

	return myKeptn, incomingEvent, err
}

func TestSyntheticCloudEventHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockCollectorIface(ctrl)

	myKeptn, incomingEvent, err := initializeTestObjects("../../test-events/collection.triggered-full.json")
	if err != nil {
		t.Error(err)
		return
	}

	eventDataHandlerIface, err := NewEventDataHandler(*incomingEvent)
	if err != nil {
		t.Errorf("Error getting keptn event data")
	}

	mockStartedEvent := cloudevents.Event{
		Context: &cloudevents.EventContextV03{},
	}
	mockStartedEvent.SetType("mock.collection.start.event")

	mockFinishedEvent := cloudevents.Event{
		Context: &cloudevents.EventContextV03{},
	}
	mockFinishedEvent.SetType("mock.collection.end.event")

	mockSyntheticTestFinishedEvent := cloudevents.Event{
		Context: &cloudevents.EventContextV03{},
	}
	mockSyntheticTestFinishedEvent.SetType("mock.synthetic.finished.event")

	m.EXPECT().GetEvents(gomock.Any()).Return([]cloudevents.Event{
		mockStartedEvent,
		mockSyntheticTestFinishedEvent,
		mockFinishedEvent,
	}, nil).AnyTimes()
	m.EXPECT().CollectExecutionIds(gomock.Any()).Return([]string{"executionId", "executionId", "executionId"}, nil)
	m.EXPECT().CollectBatchIds(gomock.Any()).Return([]string{"batchId"}, nil)
	m.EXPECT().ParseEvents(gomock.Any(), "mock.collection.start.event", "").Return([]cloudevents.Event{
		mockStartedEvent,
	})
	m.EXPECT().ParseEvents(gomock.Any(), "mock.collection.end.event", "").Return([]cloudevents.Event{
		mockFinishedEvent,
	})
	m.EXPECT().ParseEvents(gomock.Any(), "mock.synthetic.finished.event", "").Return([]cloudevents.Event{
		mockSyntheticTestFinishedEvent,
	})

	timestampA, _ := time.Parse(time.RFC3339, "2022-04-07T12:04:28Z")
	m.EXPECT().CollectEarliestTime(gomock.Any(), gomock.Any()).Return(timestampA, nil)

	timestampB, _ := time.Parse(time.RFC3339, "2022-04-07T12:05:29Z")
	m.EXPECT().CollectLatestTime(gomock.Any(), gomock.Any()).Return(timestampB, nil)

	err = CollectionCloudEventHandler(myKeptn, *incomingEvent, "serviceName", m, eventDataHandlerIface)
	assert.NilError(t, err)

	assert.Equal(t, len(myKeptn.EventSender.(*fake.EventSender).SentEvents), 2)
	assert.Equal(t, keptnv2.GetStartedEventType("collection"), myKeptn.EventSender.(*fake.EventSender).SentEvents[0].Type())
	assert.Equal(t, keptnv2.GetFinishedEventType("collection"), myKeptn.EventSender.(*fake.EventSender).SentEvents[1].Type())

	finishedEventData := CollectionSuccessfulEventData{}
	myKeptn.EventSender.(*fake.EventSender).SentEvents[1].DataAs(&finishedEventData)

	assert.Equal(t, finishedEventData.Evaluation.Start, timestampA.Format(time.RFC3339))
	assert.Equal(t, finishedEventData.Evaluation.End, timestampB.Format(time.RFC3339))

	// Test empty event
	m = NewMockCollectorIface(ctrl)

	myKeptn, incomingEvent, err = initializeTestObjects("../../test-events/collection.triggered-empty.json")
	if err != nil {
		t.Error(err)
		return
	}

	eventDataHandlerIface, err = NewEventDataHandler(*incomingEvent)
	if err != nil {
		t.Errorf("Error getting keptn event data")
	}

	m.EXPECT().GetEvents(gomock.Any()).Return([]cloudevents.Event{
		mockStartedEvent,
		mockSyntheticTestFinishedEvent,
		mockFinishedEvent,
	}, nil).AnyTimes()

	m.EXPECT().ParseEvents(gomock.Any(), "", "").Return([]cloudevents.Event{
		mockStartedEvent,
		mockFinishedEvent,
	}).Times(2)

	m.EXPECT().ParseEvents(gomock.Any(), "sh.keptn.event.test.finished", "").Return([]cloudevents.Event{
		mockSyntheticTestFinishedEvent,
	}).Times(1)

	timestampA, _ = time.Parse(time.RFC3339, "2022-04-07T12:04:28Z")
	m.EXPECT().CollectEarliestTime(gomock.Any(), gomock.Any()).Return(timestampA, nil)

	timestampB, _ = time.Parse(time.RFC3339, "2022-04-07T12:05:29Z")
	m.EXPECT().CollectLatestTime(gomock.Any(), gomock.Any()).Return(timestampB, nil)

	m.EXPECT().CollectExecutionIds(gomock.Any()).Return([]string{"executionId", "executionId", "executionId"}, nil)
	m.EXPECT().CollectBatchIds(gomock.Any()).Return([]string{"batchId"}, nil)

	err = CollectionCloudEventHandler(myKeptn, *incomingEvent, "serviceName", m, eventDataHandlerIface)
	assert.NilError(t, err)
}
