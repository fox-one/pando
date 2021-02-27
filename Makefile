.PHONY: build-server
build-server:
	sh hack/build.sh ./cmd/pando-server

.PHONY: build-worker
build-worker:
	sh hack/build.sh ./cmd/pando-worker

.PHONY: pando/worker
pando/worker:
	docker build -t pando/worker -f ./docker/Dockerfile.worker .

.PHONY: pando/server
pando/server:
	docker build -t pando/server -f ./docker/Dockerfile.server .

clean:
	@rm -rf ./builds
