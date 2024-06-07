SUBDIRS := $(wildcard cmd/*)

all: $(SUBDIRS)

$(SUBDIRS):
	$(MAKE) -C $@

clean:
	find -type f -not -path "*.git*" -executable -exec rm {} \;

.PHONY: all $(SUBDIRS) clean
