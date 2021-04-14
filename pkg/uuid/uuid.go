package uuid

import (
	uid "github.com/satori/go.uuid"
)

type GoUUIDGenerator struct {
}

var _ UUIDer = &GoUUIDGenerator{}

func NewGoUUIDGenerator() *GoUUIDGenerator {
	return &GoUUIDGenerator{}
}

func (gui *GoUUIDGenerator) V5(namespace, name string) (string, error) {
	ns, err := uid.FromString(namespace)
	if err != nil {
		return "", nil
	}
	return uid.NewV5(ns, name).String(), nil
}

func (gui *GoUUIDGenerator) V4() (string, error) {
	return uid.NewV4().String(), nil
}

func (gui *GoUUIDGenerator) V3(namespace, name string) (string, error) {
	ns, err := uid.FromString(namespace)
	if err != nil {
		return "", nil
	}
	return uid.NewV3(ns, name).String(), nil
}

func (gui *GoUUIDGenerator) V2(domain byte) (string, error) {
	return uid.NewV2(domain).String(), nil
}

func (gui *GoUUIDGenerator) V1() (string, error) {
	return uid.NewV1().String(), nil
}
