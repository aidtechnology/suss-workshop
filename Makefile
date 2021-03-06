.PHONY: all
default: help
DOCKER_IMAGE_NAME=gcr.io/fairbank-io/suss-workshop
BINARY_NAME=suss-workshop
VERSION_TAG=0.1.0

# Include build code at compile time
LD_FLAGS="\
-X github.com/aidtechnology/suss-workshop/cmd.buildCode=`git log --pretty=format:'%H' -n1`\
-X github.com/aidtechnology/suss-workshop/cmd.buildTimestamp=`date +'%s'` \
"

build: ## Build for the default architecture in use
	go build -v -ldflags $(LD_FLAGS) -o $(BINARY_NAME)

linux: ## Build for linux systems
	GOOS=linux GOARCH=amd64 go build -v -ldflags $(LD_FLAGS) -o $(BINARY_NAME)_$(VERSION_TAG)_linux

windows: ## Build for Windows systems
	GOOS=windows GOARCH=amd64 go build -v -ldflags $(LD_FLAGS) -o $(BINARY_NAME)_$(VERSION_TAG)_windows.exe

docker: ## Build docker image
	make linux
	@-docker rmi $(DOCKER_IMAGE_NAME):$(VERSION_TAG)
	@docker build --build-arg VERSION="$(VERSION_TAG)" --rm -t $(DOCKER_IMAGE_NAME):$(VERSION_TAG) .
	@rm $(BINARY_NAME)-linux

release: ## Publish the docker image to the cloud registry
	docker push $(DOCKER_IMAGE_NAME):$(VERSION_TAG)

ca-roots: ## Generate the list of valid CA certificates
	@docker run -dit --rm --name ca-roots debian:stable-slim
	@docker exec --privileged ca-roots sh -c "apt update"
	@docker exec --privileged ca-roots sh -c "apt install -y ca-certificates"
	@docker exec --privileged ca-roots sh -c "cat /etc/ssl/certs/* > /ca-roots.crt"
	@docker cp ca-roots:/ca-roots.crt ca-roots.crt
	@docker stop ca-roots

help: ## Display available make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-16s\033[0m %s\n", $$1, $$2}'
