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

func (ns *NanoIDGenerator) UUID() (string, error) {
	uuid, err := gonanoid.Nanoid(ns.Length)
	if err != nil {
		return "", err
	}
	return uuid, err
}
