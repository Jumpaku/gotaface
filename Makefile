.DEFAULT_GOAL:=help


.PHONY: help
help: ## Show this help.
	@grep -E '^[a-zA-Z_%-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%s\033[0m\n\t%s\n", $$1, $$2}'


.PHONY: lint
lint: ## lint
	find . -print | grep --regex '.*\.go' | xargs goimports -w
	go vet ./...


.PHONY: test
test: ## test
	go clean -testcache
	sh spanner/tools/integration-test.sh
	sh sqlite3/tools/integration-test.sh
	go clean -testcache
	go test -cover ./... -coverprofile=cover/cover.out && go tool cover -html=cover/cover.out -o cover/cover.html

.PHONY: test-spanner
test-spanner: ## run in container work
	SPANNER_EMULATOR_HOST=spanner:9010 \
	GOTAFACE_TEST_SPANNER_PROJECT=gotaface \
	GOTAFACE_TEST_SPANNER_INSTANCE=test \
	go test ./spanner/schema/...

.PHONY: init-spanner-test
init-spanner-test: ## initialize Spanner emulator database for develop, run in container gcloud
	gcloud config set project gotaface
	gcloud config set auth/disable_credentials true
	yes | gcloud config set api_endpoint_overrides/spanner http://spanner:9020/
	yes | gcloud spanner instances delete test || true
	gcloud spanner instances create test --config=emulator-config --description="Instance for integration test"
