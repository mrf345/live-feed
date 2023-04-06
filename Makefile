PACKAGE_NAME          := github.com/mrf345/live-feed
GOLANG_CROSS_VERSION  ?= 1.19.5
BRANCH 				  := $(shell git rev-parse --abbrev-ref HEAD)
V 					  ?= $(shell git tag)

release:
	@if [ "$(BRANCH)" != "master" ]; then\
		echo "Release must be done in master!";\
		exit 1;\
	fi
	@git pull
	@git tag v$(V)
	@docker run \
		--rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v `pwd`/sysroot:/sysroot \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--clean --skip-validate --skip-publish
	gh release create v$(V) --notes-file dist/CHANGELOG.md dist/{*.zip,checksums.txt}
	sudo rm -rf dist/ sysroot/
