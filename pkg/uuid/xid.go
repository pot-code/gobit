package uuid

import "github.com/rs/xid"

type XIDGenerator struct {
	generator xid.ID
}

var _ UUIDer = &XIDGenerator{}

func NewXIDGenerator() *XIDGenerator {
	return &XIDGenerator{xid.New()}
}

func (xg *XIDGenerator) V5(namespace, name string) (string, error) {
	return "", ErrNotImplemented
}

func (xg *XIDGenerator) V4() (string, error) {
	return xg.generator.String(), nil
}

func (xg *XIDGenerator) V3(namespace, name string) (string, error) {
	return "", ErrNotImplemented
}

func (xg *XIDGenerator) V2(domain byte) (string, error) {
	return "", ErrNotImplemented
}

func (xg *XIDGenerator) V1() (string, error) {
	return "", ErrNotImplemented
}
