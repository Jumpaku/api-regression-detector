package io

import "io"

type NamedReader interface {
	io.Reader
	Name() string
}
type NamedWriter interface {
	io.Writer
	Name() string
}
