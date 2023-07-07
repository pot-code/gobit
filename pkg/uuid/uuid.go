package uuid

import (
	uid "github.com/satori/go.uuid"
)

type GoUUIDGenerator struct {
}

var _ UUID = &GoUUIDGenerator{}

func NewGoUUIDGenerator() *GoUUIDGenerator {
	return &GoUUIDGenerator{}
}

func (gui *GoUUIDGenerator) UUID() (string, error) {
	return uid.NewV4().String(), nil
}
