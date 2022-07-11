TMPDIR := $(shell mktemp -d)
sync-protos:
	git clone --branch feat/add-query-vault-events-to-pando_gen-go git@github.com:fox-one/pando-protos.git $(TMPDIR)
	cd $(TMPDIR); git reset --hard 4f0e933236d9ebb043861155f5fe1284f0d37317
	cp -r $(TMPDIR)/pando/v1/* handler/rpc/pando
