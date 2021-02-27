.PHONY: build-server
build-server:
	sh hack/build.sh ./cmd/pando-server

.PHONY: build-worker
build-worker:
	sh hack/build.sh ./cmd/pando-worker

.PHONY: build-server
build-server:
	sh hack/build.sh ./cmd/pando-server

.PHONY: build-worker
build-worker:
	sh hack/build.sh ./cmd/pando-worker

clean:
	@rm -rf ./builds
