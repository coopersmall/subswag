CURRENT_DIR = $(shell pwd)

compile:
	@echo "Compiling..."
	@for dir in $(shell ls -d */); do \
		if [ -d "node_modules" ]; then \
			continue; \
		fi; \
		ncc build $$dir -o $(CURRENT_DIR)/$$dir/dist; \
	done
	@echo "Done!"

