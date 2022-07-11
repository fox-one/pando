TMPDIR := $(shell mktemp -d)
sync-protos:
	git clone --branch gen-go git@github.com:fox-one/pando-protos.git $(TMPDIR)
	cd $(TMPDIR); git reset --hard 3776e377a8d7112ce1dc1b830fa618b06a8595a9
	cp -r $(TMPDIR)/pando/v1/* handler/rpc/pando
