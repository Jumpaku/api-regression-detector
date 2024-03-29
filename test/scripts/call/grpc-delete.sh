#!/bin/sh

set -eux

ENDPOINT='api:50051'
METHOD='api.GreetingService/DeleteHello'
REQUEST='test/data/call/delete/request.json'
ACTUAL_RESPONSE='test/data/call/delete/actual.json'
EXPECTED_RESPONSE='test/data/call/delete/expected.json'

go run cmd/call-grpc/main.go "${ENDPOINT}" "${METHOD}" < "${REQUEST}" > "${ACTUAL_RESPONSE}"

go run cmd/compare/main.go "${EXPECTED_RESPONSE}" "${ACTUAL_RESPONSE}"