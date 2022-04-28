#!/bin/sh

docker build . -t localhost:32000/keptn-test-collector-service:dev
docker push localhost:32000/keptn-test-collector-service:dev

# Pull new image by deleting current pod
kubectl -n keptn delete pod -l app.kubernetes.io/instance=keptn-test-collector-service

# helm upgrade --install -n keptn keptn-test-collector-service chart/ --set image.repository=localhost:32000/keptn-test-collector-service --set image.tag=dev
# kubectl -n keptn logs -l app.kubernetes.io/instance=keptn-test-collector-service -c keptn-service --follow
