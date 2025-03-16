CURRENT_DIR = $(shell pwd)

compile:
	@echo "Compiling..."
	@for dir in $(shell ls -d */); do \
		if [ -f "$$dir/index.js" ]; then \
			echo "Compiling $$dir..."; \
			ncc build $$dir/index.js -o $$dir/dist; \
		fi; \
	done
	@echo "Done!"
