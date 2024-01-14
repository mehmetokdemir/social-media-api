include META

.PHONY: generate-docs
generate-docs:
	@echo "[GENERATE DOCS] Generating API documents"
	@echo " - Updating document version"
	@echo " - Initializing swag"
	@swag init --parseDependency --parseInternal --generatedTime --parseDepth 3

.PHONY: tidy
tidy:
	@echo "[TIDY] Running go mod tidy"
	@$(RUN_ENV_SET) go mod tidy -compat=1.20

.PHONY: lint
lint:
	@echo "[TIDY] Running golangci-lint run"
	@golangci-lint run

.PHONY: test
test:
	@echo "Starting unit tests"
	@go test -v  ./...


.PHONY: build
build: tidy
	@echo "[BUILD] Building the service"
	go build -o bin/social-media-api

.PHONY: run
run:
	@echo "[RUN] Running the service"
	./bin/social-media-api

.PHONY: run-with-docker
run-with-docker:
	@echo "[RUN] Running the service as a docker image"
	docker-compose up --build

.PHONY: rwgocommand
rwgocommand:
	@echo "[RUN] Running the service"
	go run .

.PHONY: git
git:
	@echo "[BUILD] Committing and pushing to remote repository"
	@echo " - Committing"
	@git add META
	@git commit -am "v$(VERSION)"
	@echo " - Tagging"
	@git tag v${VERSION}
	@echo " - Pushing"
	@git push --tags origin ${BRANCH}