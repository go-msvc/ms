package ms

import (
	"github.com/go-msvc/logger"
)

type IContext interface {
	logger.ILogger
	Cancel() func()
	GetConfig(name string) interface{}
}

func NewContext(logger logger.ILogger, configSet map[string]interface{}) IContext {
	return &context{
		ILogger: logger,
		config:  configSet,
	}
}

type context struct {
	logger.ILogger
	config map[string]interface{}
}

func (ctx context) Cancel() func() {
	return ctx.cancel
}

func (ctx context) cancel() {
}

func (ctx context) GetConfig(name string) interface{} {
	return ctx.config[name]
}
