package ms

import (
	"fmt"
	"sync"

	"github.com/go-msvc/config"
	"github.com/go-msvc/errors"
	"github.com/go-msvc/logger"
)

var (
	log = logger.ForThisPackage()
)

type IDomain interface {
	logger.ILogger

	WithOper(name string, oper IOper) IDomain
	AddOper(name string, oper IOper) error
	GetOper(name string) IOper

	NewContext(id string) IContext
}

func New(domainDefaultStructValue interface{}) IDomain {
	d := &domain{
		ILogger: log.NewLogger("domain"),
		oper:    map[string]IOper{},
	}

	if domainDefaultStructValue != nil {
		var err error
		d.cfg, err = config.Add(domainDefaultStructValue)
		if err != nil {
			panic(fmt.Sprintf("cannot create ms domain: %+v", errors.Wrapf(err, "domain config error")))
		}
	}
	return d
}

type domain struct {
	logger.ILogger
	cfg config.IConfigurable

	sync.Mutex
	oper map[string]IOper
}

func (d *domain) WithOper(name string, oper IOper) IDomain {
	if err := d.AddOper(name, oper); err != nil {
		panic(errors.Wrapf(err, "cannot add oper domain().WithOper(%s)", name))
	}
	return d
}

func (d *domain) AddOper(name string, oper IOper) error {
	d.Lock()
	defer d.Unlock()
	if _, ok := d.oper[name]; ok {
		panic(errors.Errorf("duplicate oper \"%s\" already defined.", name))
	}
	d.oper[name] = oper
	return nil
}

func (d *domain) GetOper(name string) IOper {
	d.Lock()
	defer d.Unlock()
	if oper, ok := d.oper[name]; ok {
		return oper
	}
	return UnknownOper{}
}

func (d *domain) NewContext(id string) IContext {
	return newContext(d.NewLogger(id), d.cfg)
}

type UnknownOper struct{}

func (o UnknownOper) Run(ctx IContext) (IResult, IResponse) {
	return nil, nil
}
