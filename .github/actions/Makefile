CURRENT_DIR = $(shell pwd)

compile:
	@echo "Compiling..."
	@for dir in $(shell ls -d */); do \
		if [ -d "$$dir/node_modules" ]; then \
			continue; \
		fi; \
		if [ -f "$$dir/index.js" ]; then \
			echo "Compiling $$dir..."; \
			NODE_PATH=$(CURRENT_DIR)/node_modules ncc build $$dir/index.js -o $$dir/dist; \
		fi; \
	done
	@echo "Done!"
