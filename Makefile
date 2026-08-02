APP_ENV ?= local

.PHONY: up down logs dev env cli-build cli-login cli-connect cli-open cli-release

up: env ## start web+db via docker-compose
	docker-compose up

down:
	docker-compose down

logs:
	docker-compose logs -f web

dev: env ## run web standalone (no docker) - needs Postgres reachable at DATABASE_URL
	set -a && . environments/$(APP_ENV).env && set +a && cd web && npm run dev

env: ## create environments/$(APP_ENV).env from its example, if missing
	test -f environments/$(APP_ENV).env || cp environments/$(APP_ENV).env.example environments/$(APP_ENV).env

cli-build:
	cd repotale && go build -o repotale .

cli-login: cli-build
	cd repotale && ./repotale login

cli-connect: cli-build
	cd repotale && ./repotale connect $(REPO)

cli-open: cli-build
	cd repotale && ./repotale open

# the module lives in a subdirectory, so its tags must be prefixed
# "repotale/" (e.g. `make cli-release VERSION=v0.1.1`) - a bare vX.Y.Z tag
# is read as tagging a (nonexistent) module at the repo root instead.
cli-release:
	test -n "$(VERSION)" # usage: make cli-release VERSION=v0.1.1
	git tag -a repotale/$(VERSION) -m "repotale CLI $(VERSION)"
	git push origin repotale/$(VERSION)
