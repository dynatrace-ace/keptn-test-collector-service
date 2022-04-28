package eventHandler

import (
	"fmt"
	"log"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2" // make sure to use v2 cloudevents here
	"github.com/keptn-sandbox/keptn-test-collector-service/internal/collector"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

/**
* Here are all the handler functions for the individual event
* See https://github.com/keptn/spec/blob/0.8.0-alpha/cloudevents.md for details on the payload
**/

type CollectionEventData struct {
	keptnv2.EventData
	SyntheticTestContext string `json:"syntheticTestContext"`
}

type EvaluationData struct {
	Start     string `json:"start"`
	End       string `json:"end"`
	Timeframe string `json:"timeframe"`
}

type CollectionSuccessfulEventData struct {
	keptnv2.EventData
	Evaluation EvaluationData `json:"evaluation"`
}

func CollectionCloudEventHandler(
	myKeptn *keptnv2.Keptn,
	incomingEvent cloudevents.Event,
	data *CollectionEventData,
	serviceName string,
	c collector.CollectorIface,
) error {
	log.Printf("Handling %s Event: %s", incomingEvent.Type(), incomingEvent.Context.GetID())

	// Make sure labels is not nil
	if myKeptn.Event.GetLabels() == nil {
		myKeptn.Event.SetLabels(map[string]string{})
	}

	_, err := myKeptn.SendTaskStartedEvent(data, serviceName)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to send task started CloudEvent (%s), aborting...", err.Error())
		log.Println(errMsg)
		return err
	}

	eventData := data.EventData

	syntheticExecutionEvents, _ := c.GetEvents()

	// Evaluation start is earliest test.started event timestamp
	testStartedEvents := c.ParseEventsOfType(syntheticExecutionEvents, "sh.keptn.event.test.started")
	evaluationStart, err := c.CollectEarliestTime(testStartedEvents)
	if err != nil {
		errMsg := fmt.Sprint("Failed to collect start timestamps, aborting...", err.Error())
		log.Println(errMsg)
		return err
	}

	// Evaluation end is incoming event timestamp
	evaluationEnd := incomingEvent.Context.GetTime()

	executionIds, err := c.CollectExecutionIds(syntheticExecutionEvents)
	if err != nil {
		errMsg := fmt.Sprint("Failed to collect execution ids, aborting...", err.Error())
		log.Println(errMsg)
		return err
	}

	batchIds, err := c.CollectBatchIds(syntheticExecutionEvents)
	if err != nil {
		errMsg := fmt.Sprint("Failed to collect batch ids, aborting...", err.Error())
		log.Println(errMsg)
		return err
	}

	labels := eventData.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}

	labels["executionIds"] = ""
	if len(executionIds) > 0 {
		labels["executionIds"] = strings.Join(executionIds, ",")
	}

	labels["batchIds"] = ""
	if len(batchIds) > 0 {
		labels["batchIds"] = strings.Join(batchIds, ",")
	}

	eventData.SetLabels(labels)

	successfulEventData := &CollectionSuccessfulEventData{
		EventData: eventData,
		Evaluation: EvaluationData{
			Start: evaluationStart.Format(time.RFC3339),
			End:   evaluationEnd.Format(time.RFC3339),
		},
	}

	_, err = myKeptn.SendTaskFinishedEvent(successfulEventData, serviceName)
	return err
}
