#!/bin/sh

set -eux

ENDPOINT='http://api:80/hello/[name]'
METHOD='PUT'
REQUEST='test/data/call/put/request.json'
ACTUAL_RESPONSE='test/data/call/put/actual.json'
EXPECTED_RESPONSE='test/data/call/put/expected.json'

go run cmd/call-http/main.go "${ENDPOINT}" "${METHOD}" < "${REQUEST}" > "${ACTUAL_RESPONSE}"

go run cmd/compare/main.go "${EXPECTED_RESPONSE}" "${ACTUAL_RESPONSE}"