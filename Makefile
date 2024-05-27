.DEFAULT_GOAL:=help


.PHONY: help
help: ## Show this help.
	@grep -E '^[0-9a-zA-Z_%-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%s\033[0m\n\t%s\n", $$1, $$2}'


.PHONY: lint
lint: ## lint
	find . -print | grep --regex '.*\.go' | xargs goimports -w
	go vet ./...


.PHONY: test
test: ## test
	go clean -testcache
	go test -cover ./... -coverprofile=cover/cover.out && go tool cover -html=cover/cover.out -o cover/cover.html

.PHONY: init-spanner
init-spanner: ## initialize Spanner emulator database for develop. specify INSTANCE=<instance-name>. run in container gcloud
	gcloud config set project gotaface
	gcloud config set auth/disable_credentials true
	yes | gcloud config set api_endpoint_overrides/spanner http://spanner:9020/
	SPANNER_EMULATOR_HOST=spanner:9010 \
	yes | gcloud spanner instances delete $(INSTANCE) || true
	SPANNER_EMULATOR_HOST=spanner:9010 \
	gcloud spanner instances create $(INSTANCE) --config=emulator-config --description="Instance for integration $(INSTANCE)"

.PHONY: test-spanner
test-spanner: ## run in container work
	make init-spanner INSTANCE=test
	SPANNER_EMULATOR_HOST=spanner:9010 \
	GOTAFACE_TEST_SPANNER_PROJECT=gotaface \
	GOTAFACE_TEST_SPANNER_INSTANCE=test \
	go test ./spanner/schema/... -data-source="projects/gotaface/instances/test/databases/schema"

.PHONY: example-spanner
example-spanner: ## run in container gcloud
	make init-spanner INSTANCE=example
	SPANNER_EMULATOR_HOST=spanner:9010 \
	gcloud spanner databases create --instance=example example

.PHONY: test-sqlite3
test-sqlite3: ## test-sqlite3
	GOTAFACE_TEST_SQLITE3=true go test ./sqlite3/schema/...
