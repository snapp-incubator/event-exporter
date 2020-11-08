SHELL := /bin/bash

include .makerc.dist

ifeq ($(shell test -e .makerc && echo -n yes),yes)
  include .makerc
endif

.PHONY: okd/*.yaml


okd/*.yaml:
	@echo "Processing and uploading imagestream in namespace $(NAMESPACE): ${@}" && \
	cat $@ \
		| sed "s|[<]NAMESPACE_HERE[>]|$(NAMESPACE)|g" \
		| sed "s|[<]REPOSITORY_SOURCE_SECRET_HERE[>]|$(REPOSITORY_SOURCE_SECRET)|g" \
		| sed "s|[<]OPENSHIFT_URL_HERE[>]|$(OPENSHIFT_URL_FQDN)|g" \
		| oc apply -n $(NAMESPACE) -f -


okd/:
	$(MAKE) $@*.yml
