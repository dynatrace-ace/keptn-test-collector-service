.PHONY: build

IMAGE := "dynatraceace/keptn-test-collector-service"

build: checktag
	@docker build -t "${IMAGE}:${tag}" .
	@echo "\nSuccesfully built ${IMAGE}:${tag}!"

push: build
	@docker push "${IMAGE}:${tag}"
	@echo "\nSuccesfully pushed ${IMAGE}:${tag}!"

deploy: checktag
	@helm upgrade --install -n keptn keptn-test-collector-service --set "image.tag=${tag}" chart/

checktag:
ifndef tag
$(error tag is not set)
endif
