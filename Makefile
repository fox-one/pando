TMPDIR := $(shell mktemp -d)
sync-protos:
	git clone --branch gen-go git@github.com:fox-one/pando-protos.git $(TMPDIR)
	cd $(TMPDIR); git reset --hard 1233683a902a544b68f464914ede7dbb53e08d39
	cp -r $(TMPDIR)/pando/v1/* handler/rpc/pando
