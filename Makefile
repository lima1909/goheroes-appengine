SHELL := /bin/bash
OS := $(shell uname)
GCLOUD_CMD = gcloud.cmd

ifeq ($(OS), Linux)
  GCLOUD_CMD = gcloud 
  $(info --> OS: $(OS) with CMD: $(GCLOUD_CMD))
else
  $(info --> OS: $(OS) with CMD: $(GCLOUD_CMD))
endif


deploy: 
	$(GCLOUD_CMD) app deploy -q

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

test:
	go tool vet .
	go test -race ./...