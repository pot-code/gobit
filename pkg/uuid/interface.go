package uuid

import "errors"

var (
	ErrNotImplemented = errors.New("not implemented")
)

// UUIDer UUID generator interface
type UUIDer interface {
	V1() (string, error)
	V2(domain byte) (string, error)
	V3(namespace, name string) (string, error)
	V4() (string, error)
	V5(namespace, name string) (string, error)
}
