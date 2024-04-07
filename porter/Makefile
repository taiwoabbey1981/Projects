BINDIR     := $(CURDIR)/bin
VERSION    ?= dev
SENTRY_DSN ?= 

start-dev: install setup-env-files
	bash ./scripts/dev-environment/StartDevServer.sh

run-migrate-dev: install setup-env-files
	bash ./scripts/dev-environment/RunMigrateDev.sh

install:
	bash ./scripts/dev-environment/SetupEnvironment.sh

setup-env-files:
	bash ./scripts/dev-environment/CreateDefaultEnvFiles.sh

build-cli:
	go build -ldflags="-w -s -X 'github.com/porter-dev/porter/cli/cmd/config.Version=${VERSION}' -X 'github.com/porter-dev/porter/cli/cmd/errors.SentryDSN=${SENTRY_DSN}'" -a -tags cli -o $(BINDIR)/porter ./cli

build-cli-dev:
	go build -ldflags="-X 'github.com/porter-dev/porter/cli/cmd/config.Version=${VERSION}' -X 'github.com/porter-dev/porter/cli/cmd/errors.SentryDSN=${SENTRY_DSN}'" -tags cli -o $(BINDIR)/porter ./cli

start-provisioner-dev: install setup-env-files
	bash ./scripts/dev-environment/StartProvisionerServer.sh

start-worker-dev: install setup-env-files
	bash ./scripts/dev-environment/StartWorkerServer.sh
