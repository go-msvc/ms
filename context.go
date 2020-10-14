package ms

import (
	"github.com/go-msvc/config"
	"github.com/go-msvc/logger"
)

type IContext interface {
	logger.ILogger

	Cancel() func()

	//call this in operation handlers to get
	//current copy of domainDefaultStructValue that was passed into NewDomain(...)
	//its fields are config values or constructed from config values
	DomainStruct() interface{}
}

func newContext(logger logger.ILogger, cfg config.IConfigurable) IContext {
	ctx := &context{
		ILogger: logger,
	}
	if cfg != nil {
		ctx.domainStruct, ctx.cancelStructFunc = cfg.Use()
	}
	return ctx
}

type context struct {
	logger.ILogger
	domainStruct     interface{}
	cancelStructFunc func()
}

func (ctx *context) DomainStruct() interface{} {
	return ctx.domainStruct
}

func (ctx *context) Cancel() func() {
	log.Debugf("Cancel()")
	if ctx.cancelStructFunc != nil {
		ctx.cancelStructFunc()
		ctx.cancelStructFunc = nil
	}
	return ctx.cancel
}

func (ctx context) cancel() {
	log.Debugf("cancel()")
}
