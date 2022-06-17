# Keptn Test Collector Service

### This project is based in the [keptn-service-template-go](https://github.com/keptn-sandbox/keptn-service-template-go)

## Overview

Keptn Quality Gate evaluations require test start and end timestamps (or a timespan).

However, due to the asynchronous nature of event driven architectures there is not always a definite answer to when an event actually happened. When test capabilities are furthermore distributed across multiple services timestamps become even more blurry.

The Test Collector Service is meant to be put in front of an evaluation task as part of a Keptn sequence. It's responsibility is collecting timestamps, synthetic test metadata, etc. from different Keptn contexts. Upon collection,  the Collector Service publishes a *sh.keptn.event.collection.finished* event enriched with such metadata.

## Installation

The service can be installed by cloning this repo and running:

```
helm upgrade --install -n keptn keptn-test-collector-service chart/
```

## Project setup

Example shipyard.yaml:
```
---
apiVersion: "spec.keptn.sh/0.2.0"
kind: "Shipyard"
metadata:
  name: "shipyard-simplenode"
spec:
  stages:
    - name: "dev"
      sequences:
      - name: "test"
        tasks:
        - name: "test"
      - name: "evaluation"
        tasks:
        - name: "collection"
        - name: "evaluation"
    - name: "prod"
      ...
```

## Triggering a collection

A collection can be triggered by publishing an event of type `sh.keptn.event.collection.triggered`. If no additional info is provided, the earliest and latest events in the current Keptn are parsed for test start and respectively end timestamps.

Empty event:
```
{
  "specversion": "1.0",
  "source": "<Source>",
  "type": "sh.keptn.event.collection.triggered",
  "datacontenttype": "application/json",
  "data": {
    "project": "<Keptn Project>",
    "service": "<Keptn Service>",
    "stage": "<Keptn Stage>",
    "labels": {
      "<Key>": "<Value>",
      ...
    },
    "collection": {}
  },
  "shkeptnspecversion": "0.2.4"
}
```

The following options can be provided to apply filters:

|Attribute|Required|Default|Comment|
|---|---|---|---|
|evaluationStartContext|no|Current context|Keptn context evaluation start timestamp will be parsed from. If left empty, current context will be used.|
|evaluationStartEventType|no|*|Keptn event type evaluation start timestamp will be parsed from. If left empty all events within a context will be considered.|
|evaluationEndContext|no|Current context|Keptn context evaluation end timestamp will be parsed from. If left empty, current context will be used.|
|evaluationEndEventType|no|*|Keptn event type evaluation end timestamp will be parsed from. If left empty all events within a context will be considered.|
|syntheticTestFinishedContext|no|Current context|Keptn context synthetic execution details will be parsed from. If left empty all events within a context will be considered.|
|syntheticTestFinishedEventType|no|sh.keptn.event.test.finished|Keptn event type synthetic execution details will be parsed from. If left empty all events within a context will be considered.|


Full example:
```
{
  "specversion": "1.0",
  "source": "<Source>",
  "type": "sh.keptn.event.collection.triggered",
  "datacontenttype": "application/json",
  "data": {
    "project": "<Keptn Project>",
    "service": "<Keptn Service>",
    "stage": "<Keptn Stage>",
    "labels": {
      "<Key>": "<Value>",
      ...
    },
    "collection": {
      "evaluationStartContext": "<Test start context>",
      "evaluationStartEventType": "<Test start event type>",
      "evaluationEndContext": "<Keptn context of test start time>",
      "evaluationEndEventType": "<Test end context>",
      "syntheticTestFinishedContext": "<Synthetic test finished context>",
      "syntheticTestFinishedEventType": "<Synthetic test finished event type>"
    }
  },
  "shkeptnspecversion": "0.2.4"
}
```

### A note on Synthetic test result collection

In addition to test related timestamps, the Keptn Test Collector Service also parses execution data from a synthetic test execution (more details can be found in the [Dynatrace Synthetic Service repo](https://github.com/dynatrace-ace/dynatrace-synthetic-service)).

Synthetic execution and batch ids are parsed and added as labels. A *sh.keptn.event.collection.finished* event would look similar to this:

```
{
  "data": {
    "evaluation": {
      "start": "<Test start timestamp>",
      "end": "<Test end timestamp>"
    },
    "labels": {
      "SYNTHETIC_BATCH_IDS": "<Bacth id>",
      "SYNTHETIC_EXECUTION_IDS": "<Comma seperated list of execution ids>",
      "ADDITIONAL_KEY": "<Additional labels provided in the trigger event are passed to the finished event>",
      ...
    },
    "project": "<Keptn Project>",
    "service": "<Keptn Service>",
    "stage": "<Keptn Stage>",
    "status": "succeeded"
  },
  "shkeptncontext": "<Current Keptn context>",
  "shkeptnspecversion": "0.2.4",
  "source": "keptn-test-collector-service",
  "specversion": "1.0",
  "type": "sh.keptn.event.collection.finished"
}
```

Not only is a subsequent evaluation provided with accurate timestamps, but this information can also be used to implement SLIs/SLOs as part of a Quality Gate:

sli.yaml
```
---
spec_version: '1.0'
indicators:
  e2e_page_load: "USQL;COLUMN_CHART;Page load;SELECT syntheticEvent, AVG(duration) FROM useraction WHERE usersession.internalUserId IN ($LABEL.SYNTHETIC_EXECUTION_IDS) GROUP BY syntheticEvent"
  e2e_version: "USQL;COLUMN_CHART;Version;SELECT syntheticEvent, AVG(duration) FROM useraction WHERE usersession.internalUserId IN ($LABEL.SYNTHETIC_EXECUTION_IDS) GROUP BY syntheticEvent"
  ...

```
