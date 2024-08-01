.default_goal := all

.PHONY: all
all: tidy

include scripts/make-rules/common.mk
include scripts/make-rules/golang.mk
include scripts/make-rules/tools.mk


define USAGE_OPTIONS

V 		Set verbose mode
endef
export USAGE_OPTIONS

.PHONY: tools
tools:
	@$(MAKE) tools.install

.PHONY: tidy
tidy:
	@$(GO) mod tidy