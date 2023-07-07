package uuid

type UUID interface {
	UUID() (string, error)
}
