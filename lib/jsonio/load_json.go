package jsonio

import (
	"encoding/json"
	"io"

	"github.com/Jumpaku/api-regression-detector/lib/errors"
)

type NamedReader interface {
	io.Reader
	Name() string
}

func LoadJson[T any](file io.Reader) (T, error) {
	decoder := json.NewDecoder(file)
	decoder.UseNumber()
	decoder.DisallowUnknownFields()

	var jsonValue T
	if err := decoder.Decode(&jsonValue); err != nil {
		return jsonValue, errors.Wrap(errors.BadJSON.Err(err), "fail to decode JSON")
	}

	return jsonValue, nil
}
