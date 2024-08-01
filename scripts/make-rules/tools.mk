TOOLS ?= $(BLOCKER_TOOLS) $(CRITICAL_TOOLS) $(TRIVIAL_TOOLS)

.PHONY: tools.install
tools.install: $(addprefix tools,install, $(TOOLS))

.PHONY: tools.install.%
tools.install.%:
	@echo "==========> Installing $*"
	@$(MAKE) install.$*

.PHONY: tools.verify.%
tools.verify.%:
	@echo "==========> Verifying $*"
	@$(MAKE) tools.install.$*

.PHONY: install.swagger
install.swagger:
	@$(GO) install github.com/go-swagger/go-swagger/cmd/swagger@latest

.PHONY: install.git-chglog
install.git-chglog:
	@$(GO) install github.com/git-chglog/git-chglog/cmd/git-chglog@latest

.PHONY: install.errcodegen
install.errcodegen:
	@$(GO) install $(ROOT_DIR)/tools/errcodegen/gen.go