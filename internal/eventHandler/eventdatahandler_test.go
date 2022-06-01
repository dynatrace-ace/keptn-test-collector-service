package eventHandler

import (
	"testing"

	"gotest.tools/assert"
)

func TestNewEventDataHandler(t *testing.T) {
	_, incomingEvent, err := initializeTestObjects("../../test-events/collection.triggered-full.json")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = NewEventDataHandler(*incomingEvent)
	assert.NilError(t, err)
}

func TestGetEvaluationStartContext(t *testing.T) {
	_, incomingEvent, err := initializeTestObjects("../../test-events/collection.triggered-full.json")
	if err != nil {
		t.Error(err)
		return
	}

	eventDataHandler, err := NewEventDataHandler(*incomingEvent)
	assert.NilError(t, err)

	context, err := eventDataHandler.GetEvaluationStartContext()
	assert.NilError(t, err)
	assert.Equal(t, "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", context)

	_, incomingEvent, err = initializeTestObjects("../../test-events/collection.triggered-empty.json")
	if err != nil {
		t.Error(err)
		return
	}

	eventDataHandler, err = NewEventDataHandler(*incomingEvent)
	assert.NilError(t, err)

	context, err = eventDataHandler.GetEvaluationStartContext()
	assert.NilError(t, err)
	assert.Equal(t, "0dc1538a-2550-49b5-8319-30d57a83519f", context)
}

func TestGetEvaluationEndContext(t *testing.T) {
	_, incomingEvent, err := initializeTestObjects("../../test-events/collection.triggered-full.json")
	if err != nil {
		t.Error(err)
		return
	}

	eventDataHandler, err := NewEventDataHandler(*incomingEvent)
	assert.NilError(t, err)

	context, err := eventDataHandler.GetEvaluationEndContext()
	assert.NilError(t, err)
	assert.Equal(t, "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz", context)

	_, incomingEvent, err = initializeTestObjects("../../test-events/collection.triggered-empty.json")
	if err != nil {
		t.Error(err)
		return
	}

	eventDataHandler, err = NewEventDataHandler(*incomingEvent)
	assert.NilError(t, err)

	context, err = eventDataHandler.GetEvaluationEndContext()
	assert.NilError(t, err)
	assert.Equal(t, "0dc1538a-2550-49b5-8319-30d57a83519f", context)
}
