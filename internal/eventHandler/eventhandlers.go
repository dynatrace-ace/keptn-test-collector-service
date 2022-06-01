package eventHandler

import (
	"fmt"
	"log"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2" // make sure to use v2 cloudevents here
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/keptn-sandbox/keptn-test-collector-service/internal/collector"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

/**
* Here are all the handler functions for the individual event
* See https://github.com/keptn/spec/blob/0.8.0-alpha/cloudevents.md for details on the payload
**/

// TBD:
// CollectionStartContext, if empty take current context
// CollectionStartEvent, e.g. "sh.keptn.event.test.started", if empty take earliest event

// CollectionEndContext, if empty take current context
// CollectionEndEvent, e.g. "sh.keptn.event.collection.triggered", if empty take last event

type EvaluationData struct {
	Start     string `json:"start"`
	End       string `json:"end"`
	Timeframe string `json:"timeframe"`
}

type CollectionSuccessfulEventData struct {
	keptnv2.EventData
	Evaluation EvaluationData `json:"evaluation"`
}

type CollectionUnsuccessfulEventData struct {
	keptnv2.EventData
}

func CollectionCloudEventHandler(
	myKeptn *keptnv2.Keptn,
	incomingEvent cloudevents.Event,
	// data *CollectionEventData,
	serviceName string,
	collectorIface collector.CollectorIface,
	collectionEventDataIface CollectionEventDataIface,
) error {
	log.Printf("Handling %s Event: %s", incomingEvent.Type(), incomingEvent.Context.GetID())

	// Make sure labels is not nil
	if myKeptn.Event.GetLabels() == nil {
		myKeptn.Event.SetLabels(map[string]string{})
	}

	eventData := collectionEventDataIface.GetEventData()

	_, err := myKeptn.SendTaskStartedEvent(&eventData, serviceName)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to send task started CloudEvent (%s), aborting...", err.Error())
		log.Println(errMsg)
		return err
	}

	var collectionStartEventsInContext []event.Event
	var collectionStartEventFilter string
	var collectionEndEventsInContext []event.Event
	var collectionEndEventFilter string
	var syntheticTestFinishedEventsInContext []event.Event
	var syntheticTestFinishedEventFilter string

	collectionStartContext, err := collectionEventDataIface.GetEvaluationStartContext()
	if err != nil {
		log.Println(err.Error())
		return sendTaskFail(myKeptn, eventData, serviceName, err)
	}

	collectionStartEventFilter = collectionEventDataIface.GetEvaluationStartEventFilter()

	collectionStartEventsInContext, err = collectorIface.GetEvents(collectionStartContext)
	if err != nil {
		log.Println(err.Error())
		return sendTaskFail(myKeptn, eventData, serviceName, err)
	}

	collectionEndContext, err := collectionEventDataIface.GetEvaluationEndContext()
	if err != nil {
		log.Println(err.Error())
		return sendTaskFail(myKeptn, eventData, serviceName, err)
	}

	collectionEndEventFilter = collectionEventDataIface.GetEvaluationEndEventFilter()

	isSameContextAsStart := collectionEndContext == collectionStartContext

	if isSameContextAsStart {
		collectionEndEventsInContext = collectionStartEventsInContext
	} else {
		collectionEndEventsInContext, err = collectorIface.GetEvents(collectionEndContext)
		if err != nil {
			log.Println(err.Error())
			return sendTaskFail(myKeptn, eventData, serviceName, err)
		}
	}

	isSyntheticTestFinishedContextProvided := collectionEventDataIface.IsSyntheticTestFinishedContextProvided()
	var syntheticTestFinishedContext string

	if isSyntheticTestFinishedContextProvided {
		syntheticTestFinishedContext, err = collectionEventDataIface.GetSyntheticTestFinishedContext()
		if err != nil {
			log.Println(err.Error())
			return sendTaskFail(myKeptn, eventData, serviceName, err)
		}

		syntheticTestFinishedEventFilter = collectionEventDataIface.GetSyntheticTestFinishedEventFilter()

		isSameContextAsStart = syntheticTestFinishedContext == collectionStartContext
		isSameContextAsEnd := syntheticTestFinishedContext == collectionEndContext

		if isSameContextAsStart {
			syntheticTestFinishedEventsInContext = collectionStartEventsInContext
		} else if isSameContextAsEnd {
			syntheticTestFinishedEventsInContext = collectionEndEventsInContext
		} else {
			syntheticTestFinishedEventsInContext, err = collectorIface.GetEvents(syntheticTestFinishedContext)
			if err != nil {
				log.Println(err.Error())
				return sendTaskFail(myKeptn, eventData, serviceName, err)
			}
		}
	}

	// Evaluation start is earliest event timestamp
	evaluationStartEvents := collectorIface.ParseEventsOfType(collectionStartEventsInContext, collectionStartEventFilter)

	evaluationStart, err := collectorIface.CollectEarliestTime(evaluationStartEvents)
	if err != nil {
		errMsg := fmt.Errorf("ABORTING. Failed to collect start timestamps for context %s, filtered by \"%s\": %s", collectionStartContext, collectionStartEventFilter, err.Error())
		log.Println(errMsg.Error())
		return sendTaskFail(myKeptn, eventData, serviceName, errMsg)
	}

	// Evaluation end is latest event timestamp
	evaluationEndEvents := collectorIface.ParseEventsOfType(collectionEndEventsInContext, collectionEndEventFilter)
	evaluationEnd, err := collectorIface.CollectLatestTime(evaluationEndEvents)
	if err != nil {
		errMsg := fmt.Errorf("ABORTING. Failed to collect end timestamps for context %s, filtered by \"%s\": %s", collectionEndContext, collectionEndEventFilter, err.Error())
		log.Println(errMsg.Error())
		return sendTaskFail(myKeptn, eventData, serviceName, errMsg)
	}

	if isSyntheticTestFinishedContextProvided {
		syntheticTestFinishedEvents := collectorIface.ParseEventsOfType(syntheticTestFinishedEventsInContext, syntheticTestFinishedEventFilter)

		executionIds, err := collectorIface.CollectExecutionIds(syntheticTestFinishedEvents)
		if err != nil {
			errMsg := fmt.Errorf("ABORTING. Failed to collect execution ids for context %s, filtered by \"%s\": %s", syntheticTestFinishedContext, syntheticTestFinishedEventFilter, err.Error())
			log.Println(errMsg.Error())
			return sendTaskFail(myKeptn, eventData, serviceName, errMsg)
		}

		batchIds, err := collectorIface.CollectBatchIds(syntheticTestFinishedEvents)
		if err != nil {
			errMsg := fmt.Errorf("ABORTING. Failed to collect batch ids for context %s, filtered by \"%s\": %s", syntheticTestFinishedContext, syntheticTestFinishedEventFilter, err.Error())
			log.Println(errMsg.Error())
			return sendTaskFail(myKeptn, eventData, serviceName, errMsg)
		}

		labels := eventData.GetLabels()
		if labels == nil {
			labels = map[string]string{}
		}

		labels["SYNTHETIC_EXECUTION_IDS"] = ""
		if len(executionIds) > 0 {
			labels["SYNTHETIC_EXECUTION_IDS"] = strings.Join(executionIds, ",")
		}

		labels["SYNTHETIC_BATCH_IDS"] = ""
		if len(batchIds) > 0 {
			labels["SYNTHETIC_BATCH_IDS"] = strings.Join(batchIds, ",")
		}

		eventData.SetLabels(labels)
	}

	successfulEventData := &CollectionSuccessfulEventData{
		EventData: eventData,
		Evaluation: EvaluationData{
			Start: evaluationStart.Format(time.RFC3339),
			End:   evaluationEnd.Format(time.RFC3339),
		},
	}

	return sendTaskSuccess(myKeptn, successfulEventData, serviceName)
}

func sendTaskSuccess(myKeptn *keptnv2.Keptn, data keptn.EventProperties, serviceName string) error {
	_, err := myKeptn.SendTaskFinishedEvent(data, serviceName)
	return err
}

func sendTaskFail(myKeptn *keptnv2.Keptn, eventData keptnv2.EventData, serviceName string, sourceErr error) error {
	eventData.Status = keptnv2.StatusErrored
	eventData.Result = keptnv2.ResultFailed
	eventData.Message = sourceErr.Error()

	_, err := myKeptn.SendTaskFinishedEvent(&eventData, serviceName)
	return err
}
