BINARY := webui
DOCKERVER :=`cat VERSION`
.DEFAULT_GOAL := linux

linux:
	cd cmd/$(BINARY) ; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 env go build -o $(BINARY)

docker:
	docker build  --tag="fcore/samplesearth:$(DOCKERVER)"  --file=./build/Dockerfile .

dockerlatest:
	docker build  --tag="fcore/samplesearth:latest"  --file=./build/Dockerfile .

publish: docker
	docker push fcore/samplesearth:$(DOCKERVER)
