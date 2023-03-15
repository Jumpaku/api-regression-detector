#!/bin/sh

set -eux

ENDPOINT='http://api:80/hello/[name]'
METHOD='GET'
REQUEST='test/data/call/http-get/request.json'
ACTUAL_RESPONSE='test/data/call/http-get/actual.json'
EXPECTED_RESPONSE='test/data/call/http-get/expected.json'

go run cmd/call-http/main.go "${ENDPOINT}" "${METHOD}" < "${REQUEST}" > "${ACTUAL_RESPONSE}"

go run cmd/compare/main.go "${EXPECTED_RESPONSE}" "${ACTUAL_RESPONSE}"