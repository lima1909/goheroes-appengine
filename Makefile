GCLOUD_CMD = gcloud.cmd
OS := $(shell uname)

# check the OS and user under Linux other gcloud command
checkOS:
	@if [ "$(OS)" = "Linux" ]; then \
		GCLOUD_CMD = gcloud ; \
	fi 
	@echo "--> OS: $(OS) with CMD: $(GCLOUD_CMD) "


deploy: checkOS
	$(GCLOUD_CMD) app deploy -q


version: checkOS
	$(GCLOUD_CMD) app versions list

clean: checkOS
	@for V in $(shell gcloud.cmd app versions list --format="table[no-heading](VERSION)") ; do \
        echo "delete version: $$V" ; \
		$(GCLOUD_CMD) app versions delete $$V -q ; \
    done

server:
	dev_appserver.py --enable_console --port=8082 app.yaml
	# dev_appserver.py  --clear-datastore=yes step0
	
