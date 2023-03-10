.PHONY: init-mysql
init-mysql:
	mysql --host=mysql --password=password main < examples/mysql/create.sql

.PHONY: init-postgres
init-postgres:
	psql --host=postgres --username=root --dbname=main < examples/postgres/create.sql

.PHONY: init-spanner
init-spanner:
	gcloud config set project regression-detector
	gcloud config set auth/disable_credentials true
	gcloud config set api_endpoint_overrides/spanner http://spanner:9020/
	gcloud spanner instances describe example || gcloud spanner instances create example --config=emulator-config --description="Instance for example using spanner"
	gcloud spanner databases describe main --instance=example || gcloud spanner databases create main --instance=example
	spanner-cli -p regression-detector -i example -d main --file=examples/spanner/create.sql

.PHONY: init-sqlite
init-sqlite:
	sqlite3 examples/sqlite/sqlite.db <examples/sqlite/create.sql

.PHONY: build
build:
	env GOOS=windows GOARCH=arm64 go build -ldflags '-s -w' -trimpath -o bin/windows/arm64/jrd main.go
	env GOOS=windows GOARCH=amd64 go build -ldflags '-s -w' -trimpath -o bin/windows/amd64/jrd main.go
	env GOOS=linux   GOARCH=arm64 go build -ldflags '-s -w' -trimpath -o bin/linux/arm64/jrd main.go
	env GOOS=linux   GOARCH=amd64 go build -ldflags '-s -w' -trimpath -o bin/linux/amd64/jrd main.go
	env GOOS=darwin  GOARCH=arm64 go build -ldflags '-s -w' -trimpath -o bin/darwin/arm64/jrd main.go
	env GOOS=darwin  GOARCH=amd64 go build -ldflags '-s -w' -trimpath -o bin/darwin/amd64/jrd main.go

.PHONY: lint
lint:
	find . -print | grep --regex '.*\.go' | xargs goimports -w
	golangci-lint run --fix --enable-all ./...


.PHONY: test
test:
	go test ./...
	./test/scripts/db/mysql.sh 2> /dev/null
	./test/scripts/db/postgres.sh 2> /dev/null
	./test/scripts/db/spanner.sh 2> /dev/null
	./test/scripts/db/sqlite.sh 2> /dev/null
	./test/scripts/call/grpc-failure.sh 2> /dev/null
	./test/scripts/call/grpc-success.sh 2> /dev/null
	./test/scripts/call/http-failure.sh 2> /dev/null
	./test/scripts/call/http-success.sh 2> /dev/null

.PHONY: cover
cover:
	go test -cover ./... -coverprofile=test/cover.out && go tool cover -html=test/cover.out -o test/cover.html