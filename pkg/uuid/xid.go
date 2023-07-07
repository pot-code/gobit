package uuid

import "github.com/rs/xid"

type XIDGenerator struct {
	generator xid.ID
}

var _ UUID = &XIDGenerator{}

func NewXIDGenerator() *XIDGenerator {
	return &XIDGenerator{xid.New()}
}

func (xg *XIDGenerator) UUID() (string, error) {
	return xg.generator.String(), nil
}
