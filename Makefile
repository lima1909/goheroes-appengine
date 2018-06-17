GCLOUD_CMD = gcloud.cmd

OS := $(shell uname)

ifeq ($(OS), Linux)
	GCLOUD_CMD = gcloud
endif

deploy:
	$(GCLOUD_CMD) app deploy -q


version:
	$(GCLOUD_CMD) app versions list

clean:
	@for V in $(shell gcloud.cmd app versions list --format="table[no-heading](VERSION)") ; do \
        echo "delete version: $$V" ; \
		$(GCLOUD_CMD) app versions delete $$V -q ; \
    done

server:
	dev_appserver.py --enable_console --port=8082 app.yaml
	# dev_appserver.py  --clear-datastore=yes step0
	
