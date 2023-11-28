
project = soundx
port = 8080

tag = $(project)
cname = $(project)
container-id = docker container ls --filter name='^$(cname)$$' --format "{{ .ID }}"

.PHONY: clean
clean:
	@rm -rfv build/* *~
	@go clean -cache
	@docker images -a | grep $(project) | awk '{print $3}' | xargs docker rmi

SOURCES = $(shell find . -name "*.go")

build/$(project): $(SOURCES)
	go build -o ./build/$(project) ./cmd/$(project)/main.go

build/$(project)-cli: $(SOURCES)
	go build -o ./build/$(project)-cli ./cmd/cli/main.go

.PHONY: build
build: build/$(project) build/$(project)-cli

.PHONY: test
test:
	@go test ./...

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: assert
assert:
	@for x in $(params); do test "$${!x}" || break; x=''; done; \
		if test "$${x}"; then echo "error: $${x} not set"; false; fi

dbport = 5532
dbuser = postgres
dbpass = secretsound
dbname = soundx
export DSN = postgres://$(dbuser):$(dbpass)@localhost:$(dbport)/$(dbname)?sslmode=disable

.PHONY: db
db: cname = soundx-db
db: id = $(shell docker container ls --filter name='^$(cname)$$' --format "{{ .ID }}")
db:
	@if test "$(stop)"; then \
		if test "$(id)"; then make docker-stop; docker container stop $(id); fi; \
	elif test -z "$(id)"; then id=$$(docker run \
		-it --rm --detach -p $(dbport):5432 \
		-e POSTGRES_PASSWORD=$(dbpass) \
		-e POSTGRES_USER=$(dbuser) \
		-e POSTGRES_DB=$(dbname) \
		--name $(cname) postgres:14.4); \
		if test $$? = 0; then \
			echo "$${id}"; \
			while ! docker exec -i "$${id}" pg_isready; do sleep 1; done; \
			while ! docker exec -i "$${id}" psql --quiet -U "$(dbuser)" -t --output=/dev/null -c SELECT; do sleep 1; done; \
			echo database is responding; \
		else \
			echo "error: database docker container failed to start"; \
			false; \
		fi; \
	else echo $(id); fi

.PHONY: run
run: db
	@go run cmd/$(project)/* $@

build/migrate:
	@mkdir -p $(dir $@)
	go build -tags 'postgres' -o build/ github.com/golang-migrate/migrate/v4/cmd/migrate

.PHONY: migration
migration: build/migrate
	build/migrate create -ext sql -dir db/sql $(name)

.PHONY: migrate
migrate: build/migrate
	build/migrate -path db/sql -database "$(DSN)" up

.PHONY: psql
psql:
	psql -d "$(DSN)" $(if $(sql),-c "$(sql)")

.PHONY: docker-stop
docker-stop: id = $(shell $(container-id))
docker-stop:
	@if test "$(id)"; then docker container stop $(id); fi

build/docker: Dockerfile $(SOURCES)
	docker build -t $(tag) .
	touch build/docker
	$(if $(restart-on-rebuild),make docker-stop)

.PHONY: docker-run
docker-run: dsn = $(subst @localhost:,@host.docker.internal:,$(DSN))
docker-run: restart-on-rebuild = true
docker-run: id = $(shell $(container-id))
docker-run: build/docker db
	@if test -z "$(id)"; then docker run $(if $(interactive),,--detach) -it --rm -p $(port):$(port) -e DSN="$(dsn)" --name $(cname) $(tag); else echo $(id); fi

postman-assets = api/postman.json api/postman-oapi-dl.json
$(postman-assets):
	@make assert params='key collection'
	@mkdir -p $(dir $@)
	@curl --location --silent --show-error --fail-with-body \
		--request GET 'https://api.getpostman.com/collections/$(collection)$(if $(findstring oapi,$@),/transformations)' \
		--header 'Content-Type: application/json' \
		--header 'x-api-key: $(key)' \
	> $@
	@echo "Downloaded $@ ($$(du -h $@ | cut -f1 | xargs))"

api/postman-oapi.json: api/postman-oapi-dl.json
	@mkdir -p $(dir $@)
	@jq '.output | fromjson' api/postman-oapi-dl.json > $@
	@echo "Generated $@ ($$(du -h $@ | cut -f1 | xargs))"

api/openapi.yaml: api/postman-oapi.json
	@if test ! -e $@; then yq -p=json api/postman-oapi.json > $@; echo "Generated $@ ($$(du -h $@ | cut -f1 | xargs))"; \
		else echo "Validated $@ ($$(du -h $@ | cut -f1 | xargs))"; touch $@; fi

build/oapi-codegen:
	@mkdir -p $(dir $@)
	go build -o build/ github.com/deepmap/oapi-codegen/cmd/oapi-codegen

oapi-gen-assets = api/types.gen.go api/server.gen.go api/client.gen.go api/spec.gen.go
$(oapi-gen-assets): build/oapi-codegen api/openapi.yaml
	build/oapi-codegen -package api -generate $(patsubst api/%.gen.go,%,$@) -o $@ api/openapi.yaml

.PHONY: api
api: $(oapi-gen-assets)

api-browser: cname = $(project)-swagger-ui
api-browser: port = 8081
api-browser: id = $(shell docker container ls --filter name='^$(cname)$$' --format "{{ .ID }}")
api-browser: api/openapi.yaml
	@if test -z "$(id)"; then docker run \
		--rm \
		-p $(port):8080 \
		--name "$(cname)" \
		-v "$(PWD)/api/openapi.yaml:/usr/share/nginx/html/openapi.yaml" \
		-e URLS='[{ name: "$(project)", url: "/openapi.yaml"}]' \
		-d swaggerapi/swagger-ui; fi
	@open http://localhost:$(port) || true

/usr/local/bin/postman:
	curl -o- "https://dl-cli.pstmn.io/install/osx_arm64.sh" | sh

integration: /usr/local/bin/postman api/postman.json docker-run
	postman collection run api/postman.json
