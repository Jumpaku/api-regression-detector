package cmd

import (
	"os"
	"strings"

	"github.com/Jumpaku/api-regression-detector/lib/call/http"
	libcmd "github.com/Jumpaku/api-regression-detector/lib/cmd"
	"github.com/Jumpaku/api-regression-detector/lib/cmd/cli"
	"github.com/Jumpaku/api-regression-detector/lib/errors"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio"
	"github.com/Jumpaku/api-regression-detector/lib/jsonio/wrap"
)

func RunCallHTTP(stdio *cli.Stdio, endpointURL string, method http.Method, headers []string) (code int) {
	errorInfo := errors.Info{"endpointURL": endpointURL, "method": method, "headers": headers}

	method = http.Method(strings.ToUpper(string(method)))
	switch method {
	default:
		PrintError(os.Stderr, errors.BadArgs.New(errorInfo.AppendTo("fail RunCallHTTP")))
		return 1
	case http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace:
	}

	headerMap := map[string][]string{}
	for _, header := range headers {
		key, val, err := http.ParseHeader(header)
		if err != nil {
			PrintError(os.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunCallHTTP")))
			return 1
		}

		headerMap[key] = append(headerMap[key], val)
	}

	reqBodyAny, err := jsonio.LoadJson[any](os.Stdin)
	if err != nil {
		PrintError(os.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunCallHTTP")))
		return 1
	}

	reqBody, err := wrap.FromAny(reqBodyAny)
	if err != nil {
		PrintError(os.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunCallHTTP")))
		return 1
	}

	res, err := libcmd.CallHTTP(endpointURL, method, &http.Request{Header: headerMap, Body: reqBody})
	if err != nil {
		PrintError(os.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunCallHTTP")))
		return 1
	}

	if err := jsonio.SaveJson(wrap.ToAny(res.Body), os.Stdout); err != nil {
		PrintError(os.Stderr, errors.Wrap(err, errorInfo.AppendTo("fail RunCallHTTP")))
		return 1
	}

	return 0
}
