TMPDIR := $(shell mktemp -d)
sync-protos:
	git clone --branch feat/add-pando-service_gen-go git@github.com:fox-one/pando-protos.git $(TMPDIR)
	cd $(TMPDIR); git reset --hard 8844498d7a463320eca53573409ee5e89ba64b58
	cp -r $(TMPDIR)/pando/v1 handler/rpc/pando
