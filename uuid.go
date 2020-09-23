package ms

import (
	"fmt"
	"sync"
)

func NewIdGen() IIdGen {
	return &idGen{}
}

type idGen struct {
	sync.Mutex
	next int
}

func (i *idGen) New() string {
	i.Lock()
	defer i.Unlock()
	id := fmt.Sprintf("%016x", i.next)
	i.next++
	return id
}
