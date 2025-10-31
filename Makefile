PROTO_DIR = api/proto
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)
PROTOC = protoc

go-proto:
	@echo "Generating Go gRPC files..."
	@for file in $(PROTO_FILES); do \
		echo "  -> Compilation $$file"; \
		$(PROTOC) -I=$(PROTO_DIR) $$file \
			--go_out=. \
			--go-grpc_out=. ; \
	done
	@echo "Generation complete."

YELLOW := \033[1;33m
GREEN := \033[1;32m
RESET := \033[0m

help:
	@echo "$(YELLOW)Available command:$(RESET)"
	@echo "  $(GREEN)make go-proto$(RESET)        - generating gRPC files"
.DEFAULT_GOAL := help