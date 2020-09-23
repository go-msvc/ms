package ms

import (
	"sync"

	"github.com/go-msvc/errors"
	"github.com/go-msvc/logger"
)

type IDomain interface {
	logger.ILogger

	//expose domain config but return domain to construct rest of domain
	WithConfig(name string, defaultValue interface{}) IDomain
	//get snapshot of current config at start of context
	CurrentConfig() map[string]interface{}

	WithOper(name string, oper IOper) IDomain
	AddOper(name string, oper IOper) error
	GetOper(name string) IOper
}

func NewDomain() IDomain {
	domain := &Domain{
		ILogger:    logger.NewLogger("domain"),
		IConfigSet: NewConfigSet(),
		oper:       map[string]IOper{},
	}
	return domain
}

type Domain struct {
	logger.ILogger
	IConfigSet
	sync.Mutex
	oper map[string]IOper
}

func (d *Domain) WithConfig(name string, defaultValue interface{}) IDomain {
	if _, err := d.AddConfig(name, defaultValue); err != nil {
		panic(errors.Wrapf(err, "cannot add config(%s)", name))
	}
	return d
}

func (d *Domain) CurrentConfig() map[string]interface{} {
	return d.IConfigSet.CurrentConfig()
}

func (d *Domain) WithOper(name string, oper IOper) IDomain {
	if err := d.AddOper(name, oper); err != nil {
		panic(errors.Wrapf(err, "cannot add oper domain().WithOper(%s)", name))
	}
	return d
}

func (d *Domain) AddOper(name string, oper IOper) error {
	d.Lock()
	defer d.Unlock()
	if _, ok := d.oper[name]; ok {
		panic(errors.Errorf("duplicate oper \"%s\" already defined.", name))
	}
	d.oper[name] = oper
	return nil
}

func (d *Domain) GetOper(name string) IOper {
	d.Lock()
	defer d.Unlock()
	if oper, ok := d.oper[name]; ok {
		return oper
	}
	return UnknownOper{}
}

type UnknownOper struct{}

func (o UnknownOper) Run(ctx IContext) (IResult, IResponse) {
	return nil, nil
}
