VERSION ?= $(shell git describe --tags --abbrev=0)

.PHONY: build tag release snapshot clean

build:
	@echo "🧪 Building snapshot version..."
	goreleaser release --snapshot --clean

tag:
ifndef VERSION
	$(error Please pass VERSION=vX.Y.Z)
endif
	@echo "🏷️  Tagging release as $(VERSION)"
	git tag $(VERSION)
	git push origin $(VERSION)

release: tag
	@echo "🚀 Release triggered via pushed tag $(VERSION)"

clean:
	rm -rf dist/

snapshot: build
