package eventHandler

import (
	"fmt"
	"log"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

const defaultSyntheticTestFinishedEventType = "sh.keptn.event.test.finished"

type CollectionData struct {
	EvaluationStartContext         string `json:"evaluationStartContext"`
	EvaluationStartEventType       string `json:"evaluationStartEventType"`
	EvaluationStartStage           string `json:"evaluationStartStage"`
	EvaluationEndContext           string `json:"evaluationEndContext"`
	EvaluationEndEventType         string `json:"evaluationEndEventType"`
	EvaluationEndStage             string `json:"evaluationEndStage"`
	SyntheticTestFinishedContext   string `json:"syntheticTestFinishedContext"`
	SyntheticTestFinishedEventType string `json:"syntheticTestFinishedEventType"`
	SyntheticTestFinishedStage     string `json:"syntheticTestFinishedStage"`
}

type CollectionEventData struct {
	keptnv2.EventData
	Collection   CollectionData `json:"collection"`
	eventContext cloudevents.EventContext
}

type CollectionEventDataIface interface {
	GetEventData() keptnv2.EventData
	GetEvaluationStartContext() (string, error)
	GetEvaluationStartEventFilter() string
	GetEvaluationStartStageFilter() string
	GetEvaluationEndContext() (string, error)
	GetEvaluationEndEventFilter() string
	GetEvaluationEndStageFilter() string
	GetSyntheticTestFinishedContext() (string, error)
	GetSyntheticTestFinishedEventFilter() string
	GetSyntheticTestFinishedStageFilter() string
}

/**
 * Parses a Keptn Cloud Event payload (data attribute)
 */
func parseKeptnCloudEventPayload(event cloudevents.Event, data interface{}) error {
	err := event.DataAs(data)
	if err != nil {
		log.Fatalf("Got Data Error: %s", err.Error())
		return err
	}
	return nil
}

/**
 * Parses all event data.
 */
func (collectionEventData *CollectionEventData) GetEventData() keptnv2.EventData {
	return collectionEventData.EventData
}

/**
 * Parses evaluation start context. If none was provided in event payload,
 * current context will be returned.
 */
func (collectionEventData *CollectionEventData) GetEvaluationStartContext() (string, error) {
	isProvidedByIncomingEvent := collectionEventData.Collection.EvaluationStartContext != ""
	if isProvidedByIncomingEvent {
		return collectionEventData.Collection.EvaluationStartContext, nil
	} else {
		shKeptnContextIface, err := collectionEventData.eventContext.GetExtension("shkeptncontext")
		if err != nil {
			return "", err
		}

		shKeptnContext, ok := shKeptnContextIface.(string)
		if !ok {
			return "", fmt.Errorf("error parsing Keptn context")
		}

		return shKeptnContext, nil
	}
}

/**
 * Parses evaluation start event type. If none was provided in event payload,
 * empty string filter will be returned.
 */
func (collectionEventData *CollectionEventData) GetEvaluationStartEventFilter() string {
	return collectionEventData.Collection.EvaluationStartEventType
}

/**
 * Parses evaluation start stage. If none was provided in event payload,
 * empty string filter will be returned.
 */
func (collectionEventData *CollectionEventData) GetEvaluationStartStageFilter() string {
	return collectionEventData.Collection.EvaluationStartStage
}

/**
 * Parses evaluation end context. If none was provided in event payload,
 * current context will be returned.
 */
func (collectionEventData *CollectionEventData) GetEvaluationEndContext() (string, error) {
	isProvidedByIncomingEvent := collectionEventData.Collection.EvaluationEndContext != ""

	if isProvidedByIncomingEvent {
		return collectionEventData.Collection.EvaluationEndContext, nil
	} else {
		shKeptnContextIface, err := collectionEventData.eventContext.GetExtension("shkeptncontext")
		if err != nil {
			return "", err
		}

		shKeptnContext, ok := shKeptnContextIface.(string)
		if !ok {
			return "", fmt.Errorf("error parsing Keptn context")
		}

		return shKeptnContext, nil
	}
}

/**
 * Parses evaluation end event type. If none was provided in event payload,
 * empty string filter will be returned.
 */
func (collectionEventData *CollectionEventData) GetEvaluationEndEventFilter() string {
	return collectionEventData.Collection.EvaluationEndEventType
}

/**
 * Parses evaluation end stage. If none was provided in event payload,
 * empty string filter will be returned.
 */
func (collectionEventData *CollectionEventData) GetEvaluationEndStageFilter() string {
	return collectionEventData.Collection.EvaluationEndStage
}

/**
 * Parses synthetic test finished context. If none was provided in event payload,
 * empty context will be returned.
 */
func (collectionEventData *CollectionEventData) GetSyntheticTestFinishedContext() (string, error) {
	isProvidedByIncomingEvent := collectionEventData.Collection.SyntheticTestFinishedContext != ""

	if isProvidedByIncomingEvent {
		return collectionEventData.Collection.SyntheticTestFinishedContext, nil
	} else {
		shKeptnContextIface, err := collectionEventData.eventContext.GetExtension("shkeptncontext")
		if err != nil {
			return "", err
		}

		shKeptnContext, ok := shKeptnContextIface.(string)
		if !ok {
			return "", fmt.Errorf("error parsing Keptn context")
		}

		return shKeptnContext, nil
	}
}

/**
 * Parses synthetic test finished event type. If none was provided in event payload,
 * const defaultSyntheticTestFinishedEventType will be returned.
 */
func (collectionEventData *CollectionEventData) GetSyntheticTestFinishedEventFilter() string {
	isProvidedByIncomingEvent := collectionEventData.Collection.SyntheticTestFinishedEventType != ""

	if isProvidedByIncomingEvent {
		return collectionEventData.Collection.SyntheticTestFinishedEventType
	} else {
		return defaultSyntheticTestFinishedEventType
	}
}

/**
 * Parses synthetic test finished stage. If none was provided in event payload,
 * empty string filter will be returned.
 */
func (collectionEventData *CollectionEventData) GetSyntheticTestFinishedStageFilter() string {
	return collectionEventData.Collection.SyntheticTestFinishedStage
}

func NewEventDataHandler(
	incomingEvent cloudevents.Event,
) (*CollectionEventData, error) {
	eventData := &CollectionEventData{
		eventContext: incomingEvent.Context,
	}

	err := parseKeptnCloudEventPayload(incomingEvent, eventData)
	if err != nil {
		return &CollectionEventData{}, err
	}

	return eventData, nil
}
