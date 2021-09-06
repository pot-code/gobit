package uuid

import gonanoid "github.com/matoous/go-nanoid"

// NanoIDGenerator UUID implementation using NanoID
type NanoIDGenerator struct {
	Length int
}

var _ UUID = &NanoIDGenerator{}

func NewNanoIDGenerator(length int) *NanoIDGenerator {
	if length < 1 {
		panic("length must be larger than 1")
	}
	return &NanoIDGenerator{Length: length}
}

func (ns *NanoIDGenerator) V5(namespace, name string) (string, error) {
	return "", ErrNotImplemented
}

func (ns *NanoIDGenerator) V4() (string, error) {
	uuid, err := gonanoid.Nanoid(ns.Length)
	if err != nil {
		return "", err
	}
	return uuid, err
}

func (ns *NanoIDGenerator) V3(namespace, name string) (string, error) {
	return "", ErrNotImplemented
}

func (ns *NanoIDGenerator) V2(domain byte) (string, error) {
	return "", ErrNotImplemented
}

func (ns *NanoIDGenerator) V1() (string, error) {
	return "", ErrNotImplemented
}
