SHELL := /bin/bash
OS := $(shell uname)
GCLOUD_CMD = gcloud.cmd
BRANCH := $(shell git branch 2>/dev/null  | grep '*' | sed 's/* \(.*\)/(\1)/')

ifeq ($(OS), Linux)
  GCLOUD_CMD = gcloud 
  $(info --> OS: $(OS) with CMD: $(GCLOUD_CMD))
else
  $(info --> OS: $(OS) with CMD: $(GCLOUD_CMD))
endif

branch:
	$(info --> BRANCH:  $(BRANCH))

deploy: 
	$(GCLOUD_CMD) app deploy app.yaml cron.yaml -q

logs:
	$(GCLOUD_CMD) app logs tail

ssh:
	$(GCLOUD_CMD)  alpha cloud-shell ssh

interactive:
	$(GCLOUD_CMD) alpha interactive

version: 
	@$(GCLOUD_CMD) app versions list

clean: 
	@for V in $(shell gcloud.cmd app versions list --format="table[no-heading](VERSION)") ; do \
        echo "delete version: $$V" ; \
		$(GCLOUD_CMD) app versions delete $$V -q ; \
    done

server:
	dev_appserver.cmd --enable_console --port=8082 app.yaml
	# dev_appserver.py  --clear-datastore=yes app.yaml

endpoint:
	$(GCLOUD_CMD) endpoints services list
	$(GCLOUD_CMD) endpoints services deploy swagger.yaml

prepare:
	@echo "-->" $(shell go version)
	go get -t ./...

test:
	go tool vet .
	go test -race -count=1  ./...

test-full:
	go tool vet .
	NU=TRUE go test -race -count=1  ./...
  # https://github.com/golangci/golangci-lint
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run ./...