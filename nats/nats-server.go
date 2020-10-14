package nats

import (
	"github.com/go-msvc/ms"
	"github.com/go-msvc/errors"
)

func New() ms.IServer {
	return server{}
}

type server struct {
}

func (s server) Run() error {
	return errors.Errorf("NYI")
}
