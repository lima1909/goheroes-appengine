SHELL := /bin/bash

OS := $(shell uname)
GCLOUD_CMD = gcloud.cmd

ifeq ($(OS), Linux)
    GCLOUD_CMD = gcloud 
endif


# check the OS and user under Linux other gcloud command
printCheckOS: 
	@echo "--> OS: $(OS) with CMD: $(GCLOUD_CMD) "


deploy: printCheckOS
	$(GCLOUD_CMD) app deploy -q


version: printCheckOS
	@$(GCLOUD_CMD) app versions list

clean: printCheckOS
	@for V in $(shell gcloud.cmd app versions list --format="table[no-heading](VERSION)") ; do \
        echo "delete version: $$V" ; \
		$(GCLOUD_CMD) app versions delete $$V -q ; \
    done

server:
	dev_appserver.py --enable_console --port=8082 app.yaml
	# dev_appserver.py  --clear-datastore=yes step0