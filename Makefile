TMPDIR := $(shell mktemp -d)
sync-protos:
	git clone --branch gen-go git@github.com:fox-one/pando-protos.git $(TMPDIR)
	cd $(TMPDIR); git reset --hard 9cb6d2b418a5e33a9f2d28e69874577935512b3a
	cp -r $(TMPDIR)/pando/v1/* handler/rpc/pando
