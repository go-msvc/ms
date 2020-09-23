package ms

import (
	"sync"

	"github.com/pkg/errors"
)

type IConfigSet interface {
	WithConfig(name string, defaultValue interface{}) IConfigSet
	AddConfig(name string, defaultValue interface{}) (IConfig, error)
	CurrentConfig() map[string]interface{}
}

func NewConfigSet() IConfigSet {
	return &configSet{
		configs: map[string]IConfig{},
		current: map[string]interface{}{},
	}
}

type configSet struct {
	sync.Mutex
	configs map[string]IConfig
	current map[string]interface{}
}

func (s *configSet) WithConfig(name string, defaultValue interface{}) IConfigSet {
	if _, err := s.AddConfig(name, defaultValue); err != nil {
		panic(errors.Wrapf(err, "cannot add config(%s)", name))
	}
	return s
}

func (s *configSet) AddConfig(name string, defaultValue interface{}) (IConfig, error) {
	s.Lock()
	defer s.Unlock()
	if _, ok := s.configs[name]; ok {
		return nil, errors.Errorf("configSet.Add(%s) already exists", name)
	}
	if validator, ok := defaultValue.(IValidator); ok {
		if err := validator.Validate(); err != nil {
			return nil, errors.Wrapf(err, "configSet.Add(%s): invalid default config", name)
		}
	}
	//s.configs[name] = NewConfig()
	s.current[name] = defaultValue
	return nil, nil
}

func (s *configSet) CurrentConfig() map[string]interface{} {
	s.Lock()
	defer s.Unlock()
	return s.current
}
